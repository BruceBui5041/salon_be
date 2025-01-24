package commissionmodel

type UpdateCommission struct {
	Status     *string  `json:"status"`
	Code       string   `json:"code" binding:"required"`
	RoleID     uint32   `json:"role_id" binding:"required"`
	Percentage float64  `json:"percentage" binding:"required"`
	MinAmount  *float64 `json:"min_amount"`
	MaxAmount  *float64 `json:"max_amount"`
	UpdaterID  uint32   `json:"-"`
}
