package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func main() {
	reqBody := LoginRequest{
		Username: "admin",
		Password: "123456",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Printf("Request body: %s\n", string(jsonData))

	resp, err := http.Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response body: %s\n", string(body))

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return
	}

	fmt.Printf("Login code: %d\n", loginResp.Code)
	fmt.Printf("Login message: %s\n", loginResp.Message)
	if loginResp.Code == 200 {
		fmt.Printf("Token obtained: %s\n", loginResp.Data.Token)
	}
}
