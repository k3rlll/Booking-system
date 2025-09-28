package functions

import (
	"encoding/json"
	"time"
)

type ErrDTO struct {
	Error error
	Time  time.Time
}

func (e ErrDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
