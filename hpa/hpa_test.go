package hpa

import (
	"fmt"
	"testing"
)

func TestHpa(t *testing.T) {
	s := getHistoryMetrics(TransmittedQuery, "frontend", "boutique", 60, 1)
	fmt.Println(s)
}
