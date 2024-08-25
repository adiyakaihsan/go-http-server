package types

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	CategoryID  int    `json:"category_id"`
}
