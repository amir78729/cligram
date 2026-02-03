package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var serverURL = "http://localhost:8080"

func postJSON(baseURL, endpoint string, data interface{}) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	url := baseURL + endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Success")
	} else {
		fmt.Println("Error:", resp.Status)
	}
}
