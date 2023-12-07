package internal
import (
	"time"
)

type Embedding struct {
	Embedding []float32 `json:"embedding"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"time"`
	ID        string    `json:"id"`
}
type EmbeddingRequest struct {
	Text string `json:"text"`
}
