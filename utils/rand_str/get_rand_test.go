package rand_str

import "testing"

func TestGetStrWithSymbol(t *testing.T) {
	s1 := GetStrWithSymbol(20)
	t.Log(s1)
	s2 := GetStr(20)
	t.Log(s2)
}
