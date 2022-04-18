package helper

import (
	"fmt"
	"gofuzzyclone/internal/logger"
)

// HandleError is a helper function to handle error
func HandleError(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			logger.Println("red", msg[0])
		} else {
			fmt.Println(e)
		}
		panic(e)
	}
}
