package main

import "testing"

func Test_DeleteFromRange(t *testing.T) {
	m := map[string]int{"q": 1}
	for k, _ := range m {
		delete(m, k)
	}
}
