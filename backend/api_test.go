package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestTransferFlow(t *testing.T) {
	number1 := 868505 // Sender
	number2 := 38210 // Recipient

	// Step 1: Login
	loginPayload := map[string]interface{}{
		"number":   number1,
		"password": "test1",
	}
	loginBody, _ := json.Marshal(loginPayload)

	resp, err := http.Post("http://localhost:3000/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Login failed with status: %v", resp.Status)
		t.Logf("Response body: %s", body) 
		t.FailNow()
	}

	token := resp.Header.Get("x-jwt-token")
	if token == "" {
		t.Fatal("Missing x-jwt-token in login response")
	}

	// Step 2: Transfer
	transferPayload := map[string]interface{}{
		"Recipient": number2,
		"Amount":    10,
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:3000/transfer/%d", number1), bytes.NewBuffer(transferBody))
	if err != nil {
		t.Fatalf("Creating transfer request failed: %v", err)
	}
	req.Header.Set("x-jwt-token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Transfer request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Transfer failed, status: %v", resp.Status)
	}

	// Step 3: Get Balance
	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:3000/account/%d", number1), nil)
	if err != nil {
		t.Fatalf("Creating balance request failed: %v", err)
	}
	req.Header.Set("x-jwt-token",token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Balance request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var balanceResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &balanceResp); err != nil {
		t.Fatalf("Failed to parse balance response: %v", err)
	}

	expectedBalance := float64(999990.0)
	actualBalance, ok := balanceResp["balance"].(float64)
	if !ok {
		t.Fatalf("Invalid Balance field or type")
	}

	if actualBalance != expectedBalance {
		t.Fatalf("Expected balance %.1f, got %.1f", expectedBalance, actualBalance)
	}

	t.Logf("Test passed: Balance is %.1f", actualBalance)
}
