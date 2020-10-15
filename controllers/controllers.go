package controllers

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this is where bot lives. move along")
}
