package models

type BotResponse struct {
	MessageId int `json:"update_id"`
	Message   struct {
		From struct {
			Id        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
		}
		Text string `json:"text"`
		Date int    `json:"date"`
	}
}
