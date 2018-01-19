package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bkono/go-catfindr/analyze"
)

var (
	port = flag.String("port", "", "port to listen on")
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to catfindr. Please post your txt file to /find and provide a ?min_confidence between 0.00 and 1.00"))
	return
}

func findHandler(known *analyze.Frame) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("findHandler invoked")
		minConf, err := strconv.ParseFloat(r.URL.Query().Get("min_confidence"), 64)
		if err != nil || minConf < 0.00 || minConf > 1.00 {
			http.Error(w, "min_confidence: required field, between 0.00 and 1.00, 32 bit float", http.StatusUnprocessableEntity)
			return
		}

		file, header, err := r.FormFile("frame")
		if err != nil || !strings.HasSuffix(header.Filename, ".txt") {
			http.Error(w, "frame: required file. type text/plain", http.StatusUnprocessableEntity)
			return
		}
		defer file.Close()

		// eww
		in, err := analyze.NewFrameFromReader(file) //ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "frame: failed reading due to: "+err.Error(), http.StatusInternalServerError)
			return
		}

		matches := analyze.FindMatches(in, known, minConf)

		result, err := json.Marshal(matches)
		if err != nil {
			http.Error(w, "failed encoding result to json: %s"+err.Error(), http.StatusInternalServerError)
		}
		log.Printf("evaluation done. matches: %s\n", result)

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		return
	}
}

func main() {
	flag.Parse()
	if envport := os.Getenv("PORT"); len(envport) != 0 {
		*port = envport
	}

	if len(*port) == 0 {
		log.Fatal("port is required")
	}

	f, err := os.Open("./training_image.txt")
	if err != nil {
		log.Fatalf("error opening trained image file: %v\n", err)
	}
	defer f.Close()

	known, err := analyze.NewFrameFromReader(f)
	if err != nil {
		log.Fatalf("failed converting trained image to multidimensional array: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/find", findHandler(known))

	addr := ":" + *port
	log.Printf("[server] Starting on addr: %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
