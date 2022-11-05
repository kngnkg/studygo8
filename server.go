package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
)

// 18dependency_injectionを流用
// $ ginkgo convert .
// $ ginkgo convert server_test.go

func main() {
	log.SetPrefix("[DEBUG] ")
	log.SetFlags(log.Llongfile)

	// connect to the db
	var err error
	db, err := sql.Open("postgres", "host=ch08-db user=gwp password=gwp dbname=gwp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// デバッグ用
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	// インタフェースTextをhandleRequestに渡す。
	// 構造体PostはインタフェースTextを実装しているので、handleRequestの引数にできる。
	// handleRequestは関数http.HandleFuncを返すので、返される関数はHandleFuncのメソッドシグネチャに一致する。
	// よって、最終的にURLに適したハンドラ関数を登録できる。
	//
	// sql.DBへのポインタを、構造体Postを経由して間接的にhandleRequestに渡す。
	// これによってhandleRequestに依存性を注入する。
	http.HandleFunc("/post/", handleRequest(&Post{Db: db}))

	server.ListenAndServe()
}

// main handler function
func handleRequest(t Text) http.HandlerFunc {
	// 正しいシグネチャの関数を返す
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case "GET":
			// 実際のハンドラにインタフェースTextを渡す
			err = handleGet(w, r, t)
		case "POST":
			err = handlePost(w, r, t)
		case "PUT":
			err = handlePut(w, r, t)
		case "DELETE":
			err = handleDelete(w, r, t)
		}
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Retrieve a post
// GET /post/1
func handleGet(w http.ResponseWriter, r *http.Request, post Text) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		log.Print(err)
		return
	}
	// インタフェースTextを受け入れ、構造体Postからデータを取得
	err = post.fetch(id)
	if err != nil {
		log.Print(err)
		return
	}
	output, err := json.MarshalIndent(&post, "", "\t")
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// Create a post
// POST /post/
func handlePost(w http.ResponseWriter, r *http.Request, post Text) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &post)
	err = post.create()
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(200)
	return
}

// Update a post
// PUT /post/1
func handlePut(w http.ResponseWriter, r *http.Request, post Text) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		log.Print(err)
		return
	}
	err = post.fetch(id)
	if err != nil {
		log.Print(err)
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(200)
	return
}

// Delete a post
// DELETE /post/1
func handleDelete(w http.ResponseWriter, r *http.Request, post Text) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		log.Print(err)
		return
	}
	err = post.fetch(id)
	if err != nil {
		log.Print(err)
		return
	}
	err = post.delete()
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(200)
	return
}
