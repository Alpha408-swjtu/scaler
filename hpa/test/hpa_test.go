package test

import (
	"fmt"
	"scaler/config"
	"scaler/hpa"
	"testing"
)

func TestHpa(t *testing.T) {
	h := hpa.NewHpa(config.Client, "boutique", "frontend")
	fmt.Println(h.Standard)
}

func TestExececuter(t *testing.T) {
	s := hpa.NewExecuter(config.Client, "boutique", config.Apps, 20)
	fmt.Println(s.DataMp)
}
