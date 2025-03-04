package util

import (
	"fmt"
)

func NewCheckID(name string) string {
	return fmt.Sprintf("service:%s", name)
}
