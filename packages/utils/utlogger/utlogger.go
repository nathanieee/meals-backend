package utlogger

import (
	"fmt"
	"os"
	"runtime"
)

func Error(errs ...error) {
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

func Info(data any) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fmt.Printf("\nInfo occurred at: %s:%d\nInfo: %s\n\n", frame.File, frame.Line, data)
}

func Fatal(err error) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fmt.Printf("\nFatal error occurred at: %s:%d\nError: %s\n\n", frame.File, frame.Line, err.Error())
	os.Exit(1)
}
