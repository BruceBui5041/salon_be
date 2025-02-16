package groupprovidermodel

type GroupProviderUpdate struct {
	RequesterID uint32   `json:"-"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	OwnerStrID  string   `json:"owner_id,omitempty"`
	AdminIDs    []string `json:"admin_ids,omitempty"`
}
