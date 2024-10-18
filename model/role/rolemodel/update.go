package rolemodel

type UpdateRole struct {
	Name           string            `json:"name"`
	Code           string            `json:"code"`
	Description    string            `json:"description"`
	PermissionInfo []*PermissionInfo `json:"role_permission"`
}
