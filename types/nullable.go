package types

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}

	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	}

	return nil
}

func Float64(value float64) NullFloat64 {
	return NullFloat64{
		NullFloat64: sql.NullFloat64{
			Float64: value,
			Valid:   true,
		},
	}
}
