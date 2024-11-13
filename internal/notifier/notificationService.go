package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

const (
	ContentType      = "Content-Type"
	JsonContentType  = "application/json"
	BotUrl           = "https://api.telegram.org/bot6429398075:AAFjoY4mthOBReLML8qh_-Zj_K9LZdKWQK"
	GotifyAuthHeader = "X-Gotify-Key"
)

func buildThumbnailRequest(thumbNailBytes []byte, chatId int) *http.Request {

	picUrl := BotUrl + "/sendPhoto"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	formFile, err := w.CreateFormFile("photo", "preview.jpg")

	if err != nil {
		fmt.Println("Error creating form file")
	}
	_, err = formFile.Write(thumbNailBytes)

	if err != nil {
		fmt.Println("Error writing thumbnail bytes")
	}

	ff, _ := w.CreateFormField("chat_id")
	_, err = ff.Write([]byte(fmt.Sprintf("%d", chatId)))

	if err != nil {
		fmt.Println("Error writing chatId")
	}

	err = w.Close()

	if err != nil {
		fmt.Println("Error closing formData writer")
	}

	thumbnailRequest, _ := http.NewRequest(http.MethodPost, picUrl, &b)
	thumbnailRequest.Header.Set(ContentType, w.FormDataContentType())
	return thumbnailRequest
}

func buildMessageRequest(chatId int, message string) *http.Request {
	url := BotUrl + "/sendMessage"
	a := []struct {
		ChatId int    `json:"chat_id"`
		Text   string `json:"text"`
		Photo  []byte `json:"photo"`
	}{{
		ChatId: chatId,
		Text:   message,
	}}

	messageBytes, _ := json.Marshal(a[0])
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(messageBytes))
	req.Header.Set(ContentType, JsonContentType)
	return req
}
