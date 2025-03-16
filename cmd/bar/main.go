package main

import (
	"fmt"
	"log"
	"net/http"
)

type BarHandler struct{}

func (h *BarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The bar server has received the request.")
	fmt.Fprint(w, "bar")
}

func main() {
	h := BarHandler{}

	s := http.Server{
		Addr:    ":8081",
		Handler: &h,
	}

	log.Fatal(s.ListenAndServe())
}
