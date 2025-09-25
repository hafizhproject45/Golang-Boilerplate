package utils

import "encoding/json"

type NullString struct {
	Set   bool
	Value *string
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	ns.Set = true
	if string(b) == "null" {
		ns.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.Value = &s
	return nil
}
