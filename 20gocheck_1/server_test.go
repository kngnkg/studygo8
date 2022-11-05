package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// エクスポートされた識別子はパッケージを省略してアクセスできる
	. "gopkg.in/check.v1"
)

// テストスイートの作成
type PostTestSuite struct{}

func init() {
	// テストスイートの登録
	Suite(&PostTestSuite{})
}

// パッケージtestingとの統合
func Test(t *testing.T) { TestingT(t) }

func (s *PostTestSuite) TestHandleGet(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc("/post/", handleRequest(&FakePost{}))
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/2", nil)
	mux.ServeHTTP(writer, request)

	c.Check(writer.Code, Equals, 200)
	var post Post
	json.Unmarshal(writer.Body.Bytes(), &post)
	c.Check(post.Id, Equals, 2)
}
