package errors

import "fmt"

func LogError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		return
	}
}
