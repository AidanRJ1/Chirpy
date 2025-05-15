package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
)

type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
	Error       string `json:"error"`
}

func formatJSON(data []byte) returnVals {
	var result returnVals
	json.Unmarshal(data, &result)

	return result
}

func runCurl(header string, body string, req string, url string) ([]byte, error) {
	cmd := exec.Command("curl", "-H", header, "-d", body, "-X", req, url)
	out, err := cmd.Output()
	return out, err
}

func TestFoulInput(t *testing.T) {
	body := `{"body": "I really need a kerfuffle to go to bed sooner, Fornax !"}`
	want := "I really need a **** to go to bed sooner, **** !"

	out, err := runCurl("Content-Type: application/json", body, "POST", "http://localhost:8080/api/validate_chirp")
	if err != nil {
		fmt.Printf("Error running curl cmd: %v", err)
	}

	result := formatJSON(out)

	if result.CleanedBody != want {
		t.Errorf("%s != %s", result.CleanedBody, want)
	}
}

func TestExtraInput(t *testing.T) {
	body := `{"body": "I hear Mastodon is better than Chirpy. sharbert I need to migrate", "extra": "this should be ignored"}`
	want := "I hear Mastodon is better than Chirpy. **** I need to migrate"

	out, err := runCurl("Content-Type: application/json", body, "POST", "http://localhost:8080/api/validate_chirp")
	if err != nil {
		fmt.Printf("Error running curl cmd: %v", err)
	}

	result := formatJSON(out)

	if result.CleanedBody != want {
		t.Errorf("%s != %s", result.CleanedBody, want)
	}
}

func TestCleanInput(t *testing.T) {
	body := `{"body": "I had something interesting for breakfast"}`
	want := "I had something interesting for breakfast"

	out, err := runCurl("Content-Type: application/json", body, "POST", "http://localhost:8080/api/validate_chirp")
	if err != nil {
		fmt.Printf("Error running curl cmd: %v", err)
	}

	result := formatJSON(out)

	if result.CleanedBody != want {
		t.Errorf("%s != %s", result.CleanedBody, want)
	}
}

func TestLongInput(t *testing.T) {
	body := `{"body": "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."}`
	want := "Chirp is too long"

	out, err := runCurl("Content-Type: application/json", body, "POST", "http://localhost:8080/api/validate_chirp")
	if err != nil {
		fmt.Printf("Error running curl cmd: %v", err)
	}

	result := formatJSON(out)

	if result.Error != want {
		t.Errorf("%s != %s", result.CleanedBody, want)
	}
}
