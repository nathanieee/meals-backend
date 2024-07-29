package utlogger

import (
	"encoding/json"
	"fmt"
	"log"
	"project-skbackend/packages/utils/uttelegram"
	"runtime"
)

func Error(errs ...error) {
	for _, err := range errs {
		if err != nil {
			pc := make([]uintptr, 15)
			n := runtime.Callers(2, pc)
			frames := runtime.CallersFrames(pc[:n])
			frame, _ := frames.Next()

			// * wrap the error in a custom error
			err = fmt.Errorf("Error occurred at: %s:%d, Error: %s", frame.File, frame.Line, err.Error())
			fmt.Println(err)

			// * wrap the message in the json format to be sent to telegram
			msg, _ := json.MarshalIndent(err.Error(), "", "\t")
			uttelegram.SendMessage(string(msg))
		}
	}
}

func Info(data ...any) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	jsondata, _ := json.MarshalIndent(data, "", "\t")

	fmt.Printf("\nInfo occurred at: %s:%d\nInfo: %s\n\n", frame.File, frame.Line, jsondata)
}

func Fatal(err error) {
	if err != nil {
		pc := make([]uintptr, 15)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()

		log.Fatalf("\nFatal error occurred at: %s:%d\nError: %s\n\n", frame.File, frame.Line, err.Error())
	}
}
