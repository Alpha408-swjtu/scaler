package hpa

import (
	"fmt"
	"testing"
)

func TestHpa(t *testing.T) {
	f := GetQuery("frontend", "boutique", CurrentTransmittedQuery)
	fmt.Println(f)
}
