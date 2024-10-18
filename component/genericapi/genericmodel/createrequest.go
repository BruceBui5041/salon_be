package genericmodel

type CreateRequest struct {
	Model      string      `json:"model"`
	Conditions interface{} `json:"conditions"`
	Data       interface{} `json:"data"`
}
