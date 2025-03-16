package main

import (
	"fmt"
	"log"
	"net/http"
)

type FooHandler struct{}

func (h *FooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The foo server has received the request.")
	fmt.Fprint(w, "foo")
}

func main() {
	h := FooHandler{}

	s := http.Server{
		Addr:    ":8080",
		Handler: &h,
	}

	log.Fatal(s.ListenAndServe())
}
