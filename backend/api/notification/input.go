package notification

type MarkReadInput struct {
	ID     uint64 `binding:"required" json:"id"`
	UserID uint64 `json:"user_id"`
}
