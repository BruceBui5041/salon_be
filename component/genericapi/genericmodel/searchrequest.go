package genericmodel

type SearchModelRequest struct {
	Model      string      `json:"model"`
	Conditions interface{} `json:"conditions"`
	Fields     interface{} `json:"fields"`
	OrderBy    *string     `json:"order_by,omitempty"`
}
