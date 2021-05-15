package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T){
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	hash.Add("6", "4", "2")
	// testCase
	// 2 4 6  12 14 16  22  24 26
	testCase := map[string]string{
		"2": "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k,v := range testCase{
		if hash.Get(k) != v{
			t.Errorf("Asking for %s,should have yield %s", k, v);
		}
	}

	hash.Add("8")
	// 2 4 6 8 |   12 14 16 18|  22 24 26 28
	//

	testCase["27"] = "8"

	for k,v := range testCase{
		if hash.Get(k) != v{
			t.Errorf("Asking for %s, should have yield %s", k, v);
		}
	}
}
