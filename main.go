package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/chat", chatHandler)

	e.Start(":8080")
}

func chatHandler(c echo.Context) error {
	m := new(Message)
	if err := c.Bind(m); err != nil {
		return err
	}

	response := GenerateResponseWithGPT(m.Text)

	return c.JSON(http.StatusOK, map[string]string{
		"response": response,
	})
}

type Message struct {
	Text string `json:"text"`
}

func GenerateResponseWithGPT(message string) string {
	url := "https://api.openai.com/v1/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":      "text-davinci-003", // または最新のモデルを指定
		"prompt":     message,
		"max_tokens": 50,
	})
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Request creation failed: %s", err)
	}

	req.Header.Set("Authorization", "Bearer  input KEY")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Reading response failed: %s", err)
	}

	// for dev
	log.Printf("Response: %s", body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if choices, found := result["choices"].([]interface{}); found && len(choices) > 0 {
		if firstChoice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := firstChoice["text"].(string); ok {
				return text
			}
		}
	}

	return "申し訳ありませんが、応答を生成できませんでした。"
}
