package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func TestReadReturns404(t *testing.T) {
	db, err := bolt.Open("test.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("test.db")
	defer db.Close()

	recorder := readRequest(t, db, "foo%2Fbar")

	assert.Equal(t, 404, recorder.Code)
}

func TestWrite(t *testing.T) {
	db, err := bolt.Open("test.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("test.db")
	defer db.Close()

	recorder := writeRequest(t, db, `{"url": "foo/bar", "id": 123}`)

	assert.Equal(t, 200, recorder.Code)
	assert.Empty(t, recorder.Body.String())
}

func TestReadEncodedWithMultipleDataSets(t *testing.T) {
	db, err := bolt.Open("test.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("test.db")
	defer db.Close()

	recorder := writeRequest(t, db, `{"url": "foo/bar", "id": 123}`)
	assert.Equal(t, 200, recorder.Code)
	recorder = writeRequest(t, db, `{"url": "bar/foo", "id": 124}`)
	assert.Equal(t, 200, recorder.Code)

	recorder = readRequest(t, db, "foo%2Fbar")
	assert.Equal(t, 200, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	testGetResult(t, recorder, 123, "foo/bar")

	recorder = readRequest(t, db, "bar%2Ffoo")
	assert.Equal(t, 200, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	testGetResult(t, recorder, 124, "bar/foo")
}

func TestReadNotEncoded(t *testing.T) {
	db, err := bolt.Open("test.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("test.db")
	defer db.Close()

	recorder := writeRequest(t, db, `{"url": "foo/bar", "id": 123}`)
	assert.Equal(t, 200, recorder.Code)

	recorder = readRequest(t, db, "foo/bar")
	assert.Equal(t, 200, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	testGetResult(t, recorder, 123, "foo/bar")
}

func TestFind(t *testing.T) {
	db, err := bolt.Open("test.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("test.db")
	defer db.Close()

	recorder := writeRequest(t, db, `{"url": "foo/bar", "id": 123}`)
	assert.Equal(t, 200, recorder.Code)
	recorder = writeRequest(t, db, `{"url": "foo/bar/baz", "id": 123}`)
	assert.Equal(t, 200, recorder.Code)
	recorder = writeRequest(t, db, `{"url": "bar/foo", "id": 124}`)
	assert.Equal(t, 200, recorder.Code)

	recorder = findRequest(t, db, "?id=123")

	t.Log(recorder.Body.String())

	assert.Equal(t, 200, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())

	var e [2]entry
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(&e)

	assert.Equal(t, 123, e[0].Id)
	assert.Equal(t, 123, e[1].Id)
	assert.Equal(t, "foo/bar", e[0].Url)
	assert.Equal(t, "foo/bar/baz", e[1].Url)
}

func testGetResult(t *testing.T, recorder *httptest.ResponseRecorder, id int64, url string){
	var e entry
	decoder := json.NewDecoder(recorder.Body)
	err := decoder.Decode(&e)
	assert.Nil(t, err)
	assert.Equal(t, id, e.Id)
	assert.Equal(t, url, e.Url)
}

func writeRequest(t *testing.T, db *bolt.DB, post string) *httptest.ResponseRecorder {
	// Write to db
	m := mux.NewRouter()
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "http://example.com/", strings.NewReader(post))
	assert.Nil(t, err)

	m.HandleFunc("/", addHandler(db)).Methods("POST")
	m.ServeHTTP(recorder, req)

	return recorder
}

func findRequest(t *testing.T, db *bolt.DB, query string) *httptest.ResponseRecorder {
	// Write to db
	m := mux.NewRouter()
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/" + query, nil)
	assert.Nil(t, err)

	m.HandleFunc("/", findHandler(db)).Methods("GET").Queries("id", "{id:[0-9a-zA-Z]+}")
	m.ServeHTTP(recorder, req)

	return recorder
}


func readRequest(t *testing.T, db *bolt.DB, url string) *httptest.ResponseRecorder {
	// Write to db
	m := mux.NewRouter()
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/" + url, nil)
	assert.Nil(t, err)

	m.HandleFunc("/{alias:[a-zA-Z0-9\\/]+}", getHandler(db)).Methods("GET")
	m.ServeHTTP(recorder, req)

	return recorder
}
