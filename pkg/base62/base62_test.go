package base62

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	MustInit("qJ012345679ABCDEFGHIKLM8NOPQRSTUVWXYZabdcefghijklmnoprstuvwxyz")
	m.Run()
}

func TestInt2String(t *testing.T) {
	convey.Convey("case:0", t, func() {
		seq := 0
		got := Int2String(uint64(seq))
		convey.So(got, convey.ShouldEqual, "q")
	})
	convey.Convey("case:1", t, func() {
		seq := 1
		got := Int2String(uint64(seq))
		convey.So(got, convey.ShouldEqual, "J")
	})
	convey.Convey("case:62", t, func() {
		seq := 62
		got := Int2String(uint64(seq))
		convey.So(got, convey.ShouldEqual, "Jq")
	})
	convey.Convey("case:6347", t, func() {
		seq := 6347
		got := Int2String(uint64(seq))
		convey.So(got, convey.ShouldEqual, "Jc8")
	})
}

func TestString2Int(t *testing.T) {
	convey.Convey("case:q", t, func() {
		s := "q"
		got := String2Int(s)
		convey.So(got, convey.ShouldEqual, 0)
	})
	convey.Convey("case:1", t, func() {
		s := "J"
		got := String2Int(s)
		convey.So(got, convey.ShouldEqual, 1)
	})
	convey.Convey("case:1En", t, func() {
		s := "Jc8"
		got := String2Int(s)
		convey.So(got, convey.ShouldEqual, 6347)
	})
}
