package util

import (
	"fmt"
)

func NewCheckID(name string, id string) string {
	return fmt.Sprintf("service:%s:%s", name, id)
}
