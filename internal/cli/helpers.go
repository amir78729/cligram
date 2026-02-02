package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var serverURL = "http://localhost:8080"

func postJSON(path string, data interface{}) {
	body, _ := json.Marshal(data)
	resp, err := http.Post(serverURL+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error:", resp.Status)
		return
	}
	fmt.Println("Success")
}
