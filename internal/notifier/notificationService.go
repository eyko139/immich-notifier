package notifier

import (
	"bytes"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/models"
	"mime/multipart"
	"net/http"
)

const (
	ContentType        = "Content-Type"
	JsonContentType    = "application/json"
	BotUrl             = "https://api.telegram.org/bot6429398075:AAFjoY4mthOBReLML8qh_-Zj_K9LZdKWQKc"
	GotifyAuthHeader   = "X-Gotify-Key"
	ImmichAlbumBaseUrl = "https://immich.itsmelon.com/albums/"
)

func buildThumbnailRequest(thumbNailBytes []byte, chatId int, album models.AlbumSubscription) *http.Request {

	caption := fmt.Sprintf("Update in album: <a href='%s'>%s</a>", ImmichAlbumBaseUrl+album.Id, album.AlbumName)

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

	chatIdFormField, _ := w.CreateFormField("chat_id")
	_, err = chatIdFormField.Write([]byte(fmt.Sprintf("%d", chatId)))

	captionFormField, _ := w.CreateFormField("caption")
	_, err = captionFormField.Write([]byte(caption))

	parseModeFormField, _ := w.CreateFormField("parse_mode")
	_, err = parseModeFormField.Write([]byte("HTML"))

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
