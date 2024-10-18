package rolemodel

type CreateRole struct {
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	Description string            `json:"description"`
	Permissions []*PermissionInfo `json:"role_permission"`
}

type PermissionInfo struct {
	ID               string `json:"id"`
	CreatePermission bool   `json:"create_permission"`
	ReadPermission   bool   `json:"read_permission"`
	WritePermission  bool   `json:"write_permission"`
	DeletePermission bool   `json:"delete_permission"`
}
