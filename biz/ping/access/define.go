package access

type FailPingRequest struct {
	ID   uint64 `json:"id" binding:"required"`
	Name string `json:"name"`
}
