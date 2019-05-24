package entity

type Banner struct {
	Serial int    `json:"serial"`
	Event  string `json:"event"`
	Text   string `json:"text"`
	Image  string `json:"image"`
	URL    string `json:"url"`
}
