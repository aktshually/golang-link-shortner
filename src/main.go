package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
