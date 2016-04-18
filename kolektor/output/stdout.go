
package output

import (
  "fmt"
)

type StdoutBackend struct {
}

func NewStdoutBackend(beArgs map[string]string) (*StdoutBackend,error) {
	return &StdoutBackend{},nil
}

func (b *StdoutBackend) Flush(metrics []Metric) {
  fmt.Println(metrics)
}
