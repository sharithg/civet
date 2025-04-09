package outing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
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

type Friend struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Outing struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	Status        string    `json:"status"`
	Friends       []Friend  `json:"friends"`
	TotalReceipts int64     `json:"total_receipts"`
}

func toOutingsResponse(outings []repository.GetOutingsRow) ([]Outing, error) {
	var outingsResp []Outing

	for _, outing := range outings {
		var friends []Friend
		if err := json.Unmarshal(outing.Friends, &friends); err != nil {
			log.Printf("error decoding friends JSON: %v", err)
			return nil, err
		}

		outingsResp = append(outingsResp, Outing{
			ID:            outing.ID,
			Name:          outing.Name,
			CreatedAt:     outing.CreatedAt.Time,
			Status:        outing.Status,
			TotalReceipts: outing.TotalReceipts,
			Friends:       friends,
		})
	}

	return outingsResp, nil
}
