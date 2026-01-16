package recoverx

import (
	"fmt"
	"os"
	"runtime/debug"
)

func CatchPanicAndDebugPrint() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", r)
		debug.PrintStack()
	}
}
