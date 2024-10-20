package categorymodel

type UpdateCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}
