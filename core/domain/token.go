package domain

type Token struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}
