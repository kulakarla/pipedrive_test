package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	apiToken      = "863be942d8456f146e61026f7cf69dc78efda801"
	apiTokenParam = "?api_token=" + apiToken
	baseUrl       = "https://api.pipedrive.com/v1/deals/"
)

func main() {
	http.HandleFunc("/deals", handler)
	http.HandleFunc("/deals/", handler)
	log.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyRequest(rw, req)
	default:
		http.Error(rw, "Method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
	}
}

func proxyRequest(rw http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	path := strings.TrimPrefix(req.URL.Path, "/deals")
	targetURL := baseUrl + path + apiTokenParam

	proxyReq, err := http.NewRequest(req.Method, targetURL, req.Body)
	if err != nil {
		http.Error(rw, "Failed to create request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = req.Header.Clone()

	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(rw, "Failed to get response from target URL", http.StatusInternalServerError)
		return
	}

	defer proxyResp.Body.Close()

	for name, values := range proxyResp.Header {
		for _, value := range values {
			rw.Header().Add(name, value)
		}
	}

	rw.WriteHeader(proxyResp.StatusCode)

	_, err = io.Copy(rw, proxyResp.Body)
	if err != nil {
		http.Error(rw, "Failed to copy response body", http.StatusInternalServerError)
	}

}
