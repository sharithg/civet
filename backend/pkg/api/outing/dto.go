package outing

import (
	"github.com/sharithg/civet/internal/repository"
	"github.com/sharithg/civet/pkg/api/utils"
)

type OutingData struct {
	id            string
	name          string
	totalReceipts int64
	createdAt     string
	friends       int64
	status        string
}

type GetReceipt struct {
	Restaurant string   `json:"restaurant"`
	OrderCount int64    `json:"order_count"`
	Total      *float64 `json:"total"`
	ID         string   `json:"id"`
}

func toOutingReceiptsResponse(receipts []repository.GetReceiptsForOutingRow) []GetReceipt {
	var rec []GetReceipt

	for _, r := range receipts {
		rec = append(rec, GetReceipt{
			Restaurant: r.Restaurant,
			OrderCount: r.OrderCount,
			Total:      utils.NullFloat64ToPtr(r.Total),
			ID:         r.ID.String(),
		})
	}

	return rec
}
