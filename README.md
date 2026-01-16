# Scalable Go Unit Test Structure

This isn't "yet another tutorial" about why unit testing is good or what Table-Driven Testing is. This is about **how you structure your unit test code efficiently**.

## Motivation

Way before my Go career, I relied heavily on unit testing in C/C++. I picked up good habits, but I always faced a painful issue: unscalable, complex tests with too much boilerplate.

When I turned to Go, I faced the same issues. I couldn't find a structure online that solved these problems well. So, I brewed my own solution.

**This structure scales well, reduces complexity, and minimizes boilerplate.** It is designed to be a copy-paste solution.

## Solution

### 1. General Rules

* **Dependency Injection:** Always use DI. We mock objects because we don't want to test dependencies of the unit.
* **Table-Driven Testing:** Always use it. Even if you have few test cases now, you don't know when you will need more.
* **Black Box Testing:** Test only the exported functions. Whitebox testing locks the internal structure of the unit and makes refactoring hard.
* **Copy-Paste:** Use the provided files as a template. Not much needs to change to fit your code.
* **Use `testify` for asserting/requiring and `gomock` for mocking:** I have done a comparison among the available frameworks. I found the combination of `testify` and `gomock` gives the best control over setting assertions and expectations, and they log in an understandable way.

### 2. File Structure

The unit test related files shall be in the same folder as the unit.

* `iam.go`: The unit itself.
* `iam_test.go`: The unit tests. (Note: unit name with `_test.go` suffix).
* `iam_main_test.go`: Contains common functions (bootstrap, main, mock gen). This separates "plumbing" from actual tests. (Note: package name with `_main_test.go` suffix).
* `mock/`: A folder next to the unit.
* `mock/iam_mock.go`: Contains all generated mocks for this package. (Note: Note: package name with `_mock.go` suffix).

### 3. The Main Test File (`*_main_test.go`)

This file holds the logic that makes the test files clean.

**Mock Generation:** Only this file contains the directive. It makes it clear where mocks come from. It's good habit to generate all mocks in one file.
```go
//go:generate go run go.uber.org/mock/mockgen -destination=mock/iam_mock.go -package=mock . identityProvider // Change: destination, interface(s) only
```

**Mocks Struct:** Group all mocks in one struct. It scales nicely; if you add a dependency, you change it here, and it is available in boostrap and tests.
```go
type mocks struct {
	idp *mock.MockidentityProvider // Change: only the field(s) of struct
}
```

**Test Main:** Always check wether any goroutine is leaking. Even if the current implementation of unit does not use goroutine, this can change anytime.
```go
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
```

**Bootstrap Function:** This function does everything to prepare the test. It creates the controller, the mocks, and the System Under Test (SUT). This is a reusable function across the specific package. You just need to change the `mocks` and how to `new` the SUT.
```go
func bootstrap[TArg any, TWant any](
	t *testing.T,
	arg TArg,
	want TWant,
	prepare func(ctx context.Context, m *mocks, param TArg, want TWant),
) (context.Context, *iam.IAM) {
	defer recoverx.CatchPanicAndDebugPrint()

	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	ctx := t.Context()

	m := mocks{
		idp: mock.NewMockidentityProvider(ctrl), // Change: how to new mock(s)
	}

	if prepare != nil {
		prepare(ctx, &m, arg, want)
	}

	sut := iam.NewIAM(m.idp) // Change: how to new unit

	return ctx, sut
}
```

**Internal Dependency Helpers:** If a function calls another exported function of the same unit (e.g., `UserByEmail` calls `ListUsers`), define helpers here to standardize the expectation.
```go
func ListUsers_Succeeds(paramCtx context.Context, wantUsers []*user.User, idp *mock.MockidentityProvider) {
  idp.EXPECT().ListUsers(paramCtx).Return(wantUsers, nil)
}

func ListUsers_Fails(paramCtx context.Context, wantErr error, idp *mock.MockidentityProvider) {
	idp.EXPECT().ListUsers(paramCtx).Return(nil, wantErr)
}
```

### 4. The Test File (`*_test.go`)

* **Parallelism:** Use `t.Parallel()` in the function and in each test case.
* **Structs:** Define `arg` and `want` structs locally in each test. Obviously, `arg` contains the passing values to the function, `want` contains the expected return values of the function.
* **Error Handling:** The `want.err` should be `any`. This allows checking for `error` (Is), `string` (Contains), or `nil` flexibly using `testx.AssertError`.

**The Skeleton:** All tests shall follow this skeleton:
```go
func TestUserByID(t *testing.T) {
	t.Parallel()

	type arg struct {
		userID string // Change: field(s) based on the parameter(s) of the function
	}

	type want struct {
		user *user.User // Change: field(s) based on the return values of the function
		err  any
	}

	tcs := []struct {
		name    string
		arg     arg
		prepare func(ctx context.Context, m *mocks, arg arg, want want)
		want    want
	}{
		{ // Change: obviously write your own TC(s) here
			name: "user has found",
			arg: arg{
				userID: "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
			},
			prepare: func(ctx context.Context, m *mocks, arg arg, want want) {
				m.idp.EXPECT().FetchUser(ctx, arg.userID).Return(want.user, nil)
			},
			want: want{
				user: &user.User{
					ID:       "62f36ff8-5ac9-4b2d-8049-197bdea5a48b",
					Username: "John Doe",
					Email:    "john@doe.com",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, sut := bootstrap(t, tc.arg, tc.want, tc.prepare)

            // Change: the rest based on function and assertion(s)
			gotUser, gotErr := sut.UserByID(ctx, tc.arg.userID)

			assert.Equal(t, tc.want.user, gotUser, tc.name)
			testx.AssertError(t, tc.want.err, gotErr, tc.name)
		})
	}
}
```

### 5. Misc

#### 5.1 Matcher for gomock

Matchers shall live in Main Test File.

Strongly recommended to add `recovery` to `Matches()` because a `panic` happens the mocking framework will just swallow and there won't be any useful log about the panic.

```go
type userMatcher struct {
	want user.User
}

func (m userMatcher) Matches(x any) bool {
	defer recoverx.CatchPanicAndDebugPrint()

	got, ok := x.(*user.User)
	if !ok {
		return false
	}

	diff := cmp.Diff(m.want, *got)
	if diff != "" {
		return false
	}

	return true
}

// This shows what we WANTED
func (m userMatcher) String() string {
	defer recoverx.CatchPanicAndDebugPrint()

	return fmt.Sprintf("%+v", m.want)
}

// This shows what we actually GOT
func (m userMatcher) Got(x any) string {
	defer recoverx.CatchPanicAndDebugPrint()

	got, ok := x.(*user.User)
	if !ok {
		return fmt.Sprintf("is not a *user.User (type %T)", x)
	}

	if got == nil {
		return "nil"
	}

	return fmt.Sprintf("%+v", *got)
}

func UserMatches(u user.User) gomock.Matcher {
	return userMatcher{want: u}
}

...

m.db.EXPECT().User(ctx, UserMatches(user.User{...})).Return(...)
```
