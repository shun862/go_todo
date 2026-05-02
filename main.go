package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello TODO App")
}
