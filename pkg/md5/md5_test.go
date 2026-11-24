package md5

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	convey.Convey("基础示例", t, func() {
		data := []byte("hhhhhhhhhlllllllloooooooo")
		got := Sum(data)
		convey.So(got, convey.ShouldEqual, "78ad9ea5e7db04e2b87e9514f98815ae")
	})

	convey.Convey("空输入", t, func() {
		data := []byte("")
		got := Sum(data)
		convey.So(got, convey.ShouldEqual, "d41d8cd98f00b204e9800998ecf8427e")
	})
}
