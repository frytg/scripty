package main

import (
	"os"
	"bytes"
	"context"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type TokenRequest struct {
	Aud string `json:"aud"`
}

type TokenResponse struct {
	Version    int    `json:"version"`
	Success    bool   `json:"success"`
	TokenType  string `json:"token_type"`
	IDToken    string `json:"id_token"`
	Expiration int64  `json:"expiration_time"`
}

// Add this struct to store the JWT payload claims
type JWTPayload struct {
	Exp int64 `json:"exp"`
	// Add other claims you need here
}

func main() {
	// Create Unix domain socket transport
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", "/.fly/api")
		},
	}
	client := &http.Client{Transport: transport}

	// Prepare request body
	reqBody := TokenRequest{
		Aud: os.Args[1], // use audience from command line
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", "http://localhost/v1/tokens/oidc", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Errorf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response and trim any whitespace including newlines
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error reading response: %v\n", err)
		return
	}
	token := strings.TrimSpace(string(body))

	// If you need to parse the JWT token parts:
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Errorf("invalid JWT token format")
		return
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Printf("error decoding JWT payload: %v\n", err)
		return
	}

	// Parse the JWT payload JSON
	var claims JWTPayload
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		fmt.Printf("error parsing JWT claims: %v\n", err)
		return
	}

	// Create formatted response
	output := TokenResponse{
		Version:    1,
		Success:    true,
		TokenType:  "urn:ietf:params:oauth:token-type:id_token",
		IDToken:    token, // full unparsed token
		Expiration: claims.Exp, // use expiration from JWT claims
	}

	// Print JSON response
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting output: %v\n", err)
		return
	}
	fmt.Println(string(jsonOutput))
}