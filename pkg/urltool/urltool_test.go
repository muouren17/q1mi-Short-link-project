package urltool_test

import (
	"shortener/pkg/urltool"
	"testing"
)

func TestGetBasePath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targeturl string
		want      string
		wantErr   bool
	}{
		{name: "基本测试", targeturl: "https://liwenzhou.com/courses/go-micro-service/short-url/#%e9%9c%80%e6%b1%82%e6%8f%8f%e8%bf%b0", want: "short-url", wantErr: false},
		{name: "无效url示例", targeturl: "/1242113/.com", want: "", wantErr: true},
		{name: "空", targeturl: "", want: "", wantErr: true},
		{name: "带querystring的链接", targeturl: "https://www.bing.com/search?q=gtihub&form=ANNTH1", want: "search", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := urltool.GetBasePath(tt.targeturl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBasePath() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBasePath() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("GetBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
