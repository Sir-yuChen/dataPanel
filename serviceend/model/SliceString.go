package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SliceString []string

func (a *SliceString) Scan(src any) error {
	jsonB, ok := src.([]byte)
	if !ok {
		return errors.New("source is not a byte array")
	}
	if !json.Valid(jsonB) {
		return errors.New("invalid json data")
	}
	return json.Unmarshal(jsonB, a)
}

func (a SliceString) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	jStr, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return []byte(jStr), nil
}
