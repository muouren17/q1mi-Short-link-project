package connect

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	convey.Convey("基础用例", t, func() {
		url := "https://liwenzhou.com/courses/go-micro-service/short-url/#%e9%9c%80%e6%b1%82%e6%8f%8f%e8%bf%b0"
		got := Get(url)
		convey.ShouldBeTrue(got)
	})

	convey.Convey("请求不通", t, func() {
		url := "123./235"
		got := Get(url)
		convey.So(got, convey.ShouldEqual, false)
		// convey.So(got, convey.ShouldBeTrue) // 不能这样写
	})
}
