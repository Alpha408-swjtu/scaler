package hpa

import (
	"fmt"
	"testing"
)

func TestHpa(t *testing.T) {
	s := GetQps("frontend", "boutique", TransmittedQuery)

	fmt.Println(s)
}
