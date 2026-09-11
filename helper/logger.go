package helper

import (
	"log"
)

func Logger(httpMethod string, statusCode int, message string) {
	log.Printf("%s /products failed: statusCode=%d message=%s", httpMethod, statusCode, message)
}
