package object

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
	tests := []string{
		"Hello World",
		"My name is johnny",
		"",
	}

	for _, str := range tests {
		str1 := &String{Value: str}
		str2 := &String{Value: str}

		if str1.HashKey() != str2.HashKey() {
			t.Errorf("strings with same content have different hash keys")
		}
	}

	for i := 0; i < len(tests)-1; i++ {
		str1 := &String{Value: tests[i]}
		for j := i + 1; j < len(tests); j++ {
			str2 := &String{Value: tests[j]}
			if str1.HashKey() == str2.HashKey() {
				t.Errorf("strings with different content have same hash keys")
			}
		}
	}
}
