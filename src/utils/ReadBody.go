package utils

import (
	"io"
	"log"
)

func ReadBody(body io.ReadCloser) []byte {
	defer body.Close()

	response, err := io.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
