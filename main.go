package main

import (
    "fmt"
    "net/http"
	"io/ioutil"
)


func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello World!")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
	}
	w.Write(body)
}

func main() {
    fmt.Println("Hello, Go!")
	fmt.Println("Starting server on port 8080")

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server")
	}
}