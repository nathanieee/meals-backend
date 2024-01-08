package utlogger

import (
	"fmt"
	"runtime"
)

func LogError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			pc := make([]uintptr, 15)
			n := runtime.Callers(2, pc)
			frames := runtime.CallersFrames(pc[:n])
			frame, _ := frames.Next()

			fmt.Printf("\nError occurred at: %s:%d\nError: %s\n\n", frame.File, frame.Line, err.Error())
		}
	}
}
