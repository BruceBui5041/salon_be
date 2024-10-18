package permissionmodel

type CreatePermission struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}
