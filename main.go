package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type UserAttributes struct {
	Attr1 string `json:"attr1"`
	Attr2 string `json:"attr2"`
}

type Payload struct {
	UserID               string         `json:"userId"`
	EntityAttributes     map[string]int `json:"entityAttributes"`
	DecideOptions        []string       `json:"decideOptions"`
	FetchSegmentsOptions []string       `json:"fetchSegmentsOptions"`
	DecisionFlagKeys     []string       `json:"decisionFlagKeys"`
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func main() {
	client := &http.Client{}
	data := Payload{
		UserID: "123",
		EntityAttributes: map[string]int{
			"interest": 10,
		},
		DecideOptions: []string{
			"DISABLE_DECISION_EVENT",
		},
		FetchSegmentsOptions: []string{
			"IGNORE_CACHE",
		},
		DecisionFlagKeys: []string{
			"opti_e2e_test",
		},
	}
	for {
		var err error
		data.UserID, err = generateRandomString(8)
		if err != nil {
			log.Fatal(err)
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", os.Getenv("DECIDE_URL"), bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Optimizely-SDK-Key", os.Getenv("X-Optimizely-SDK-Key"))

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatal("Failed to get response with status code: ", resp.StatusCode)
		}

		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// resp.Body.Close()

		// fmt.Println("Response body: ", string(bodyBytes))
		// fmt.Println("Status Code: ", resp.StatusCode)
		resp.Body.Close()
		resp.Close = true
		req.Close = true

		time.Sleep(100 * time.Millisecond)
	}
}
