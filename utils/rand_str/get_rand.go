package rand_str

import (
	"math/rand"
	"time"
)

const (
	ALLCHARS    = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	ALLWORDCHAR = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GetStrWithSymbol(l int) string {
	rand.Seed(time.Now().Unix())
	randChars := make([]byte, 0)
	rang := len(ALLCHARS)
	for i := 0; i < l; i++ {
		randChars = append(randChars, ALLCHARS[rand.Intn(rang)])
	}
	return string(randChars)
}

func GetStr(l int) string {
	rand.Seed(time.Now().UnixNano())
	randChars := make([]byte, 0)
	rang := len(ALLWORDCHAR)
	for i := 0; i < l; i++ {
		randChars = append(randChars, ALLWORDCHAR[rand.Intn(rang)])
	}
	return string(randChars)
}
