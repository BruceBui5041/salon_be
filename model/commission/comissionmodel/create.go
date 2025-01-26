package commissionmodel

type CreateCommission struct {
	Code       string   `json:"code" binding:"required"`
	RoleIDStr  string   `json:"role_id" binding:"required"`
	RoleID     uint32   `json:"-"`
	Percentage float64  `json:"percentage" binding:"required"`
	MinAmount  *float64 `json:"min_amount"`
	MaxAmount  *float64 `json:"max_amount"`
	CreatorID  uint32   `json:"-"`
}
