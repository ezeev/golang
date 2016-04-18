package output

import (
	"errors"
)

type Metric struct {
  Name string
  Value float64
  Timestamp int64
  Tags map[string]string
}



type Backend interface {
	Flush([]Metric)
}

func NewBackend(be string, beArgs map[string]string) (Backend,error) {
	if (be == "wavefront") {
		return NewWavefrontBackend(beArgs)
	} else if (be == "stdout") {
		return NewStdoutBackend(beArgs)
	} else {
		return nil, errors.New("No or invalid backend specified!")
	}
}
