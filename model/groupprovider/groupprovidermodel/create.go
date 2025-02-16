package groupprovidermodel

type GroupProviderCreate struct {
	RequesterID uint32 `json:"-"`
	OwnerStrID  string `json:"owner_id"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}
