package pkg

import "log"

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Println(message+": %v", err.Error())
	}
}