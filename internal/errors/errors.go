package errors

import "log"

func HandleErrorWithMessage(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
