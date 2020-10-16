package main

import (
	"fmt"
	"log"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kondition is live and running. Navigate to a path to check the status of a watched service.")
}

func main() {
	http.HandleFunc("/", defaultHandler)
	log.Print("Kondition is live and listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
