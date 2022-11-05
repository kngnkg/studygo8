package main

// go test -check.vv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// エクスポートされた識別子はパッケージを省略してアクセスできる
	. "gopkg.in/check.v1"
)

// テストスイートの作成
// テストスイートに保存されるテストフィクスチャのデータ
type PostTestSuite struct {
	mux    *http.ServeMux
	post   *FakePost
	writer *httptest.ResponseRecorder
}

func init() {
	// テストスイートの登録
	Suite(&PostTestSuite{})
}

// パッケージtestingとの統合
func Test(t *testing.T) { TestingT(t) }

func (s *PostTestSuite) SetUpTest(c *C) {
	s.post = &FakePost{}
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/post/", handleRequest(s.post))
	s.writer = httptest.NewRecorder()
}

func (s *PostTestSuite) TearDownTest(c *C) {
	c.Log("Finished test - ", c.TestName())
}

func (s *PostTestSuite) SetUpSuite(c *C) {
	c.Log("Starting Post Test Suite")
}

func (s *PostTestSuite) TearDownSuite(c *C) {
	c.Log("Finishing Post Test Suite")
}

func (s *PostTestSuite) TestHandleGet(c *C) {
	request, _ := http.NewRequest("GET", "/post/2", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	var post Post
	json.Unmarshal(s.writer.Body.Bytes(), &post)
	c.Check(post.Id, Equals, 2)
}

func (s *PostTestSuite) TestHandlePut(c *C) {
	json := strings.NewReader(`{"content":"Updated post","author":"Sau Sheong"}`)
	request, _ := http.NewRequest("PUT", "/post/2", json)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	c.Check(s.post.Id, Equals, 2)
	c.Check(s.post.Content, Equals, "Updated post")
}
