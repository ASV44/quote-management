// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          int32           `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Price       float64         `json:"price"`
	TaxRate     float64         `json:"taxRate"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
