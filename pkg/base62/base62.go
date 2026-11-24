package base62

import (
	"math"
	"slices"
	"strings"
)

var base62Str string

func MustInit(bs string) {
	if len(bs) == 0 {
		panic("must need baseStr!")
	}
	base62Str = bs
}

func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	res := make([]string, 0)
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		res = append(res, string(base62Str[mod]))
		seq = div
	}
	slices.Reverse(res)
	return strings.Join(res, "")
}

func String2Int(s string) uint64 {
	bl := []byte(s)
	slices.Reverse(bl)
	res := 0.0
	for i, b := range bl {
		idx := strings.IndexByte(base62Str, b)
		res += math.Pow(62, float64(i)) * float64(idx)
	}
	return uint64(res)
}
