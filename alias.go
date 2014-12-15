package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/boltdb/bolt"
	"net/http"
	"flag"
	"encoding/json"
	"encoding/binary"
	"fmt"
)

type entry struct {
	Url string 		`json:"url"`
	Id int64		`json:"id"`
	Created int64	`json:"created"`
	Updated int64	`json:"updated"`
}

func main() {
	portPtr := flag.String("port", "80", "The port to listen on")
	hostPtr := flag.String("host", "", "The adress to listen on. Leave empty to listen on all interfaces")
	dbLocPtr := flag.String("db", "alias.db", "A path pointing to the database.")
	flag.Parse()

	listen := fmt.Sprintf("%s:%s", *hostPtr, *portPtr)

	db, err := bolt.Open(*dbLocPtr, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/{alias:[a-zA-Z0-9\\/]+}", getHandler(db)).Methods("GET")
	r.HandleFunc("/", findHandler(db)).Methods("GET").Queries("id", "{id:[0-9a-zA-Z]+}")
	r.HandleFunc("/", addHandler(db)).Methods("POST")

	log.Fatal(http.ListenAndServe(listen, r))
}

func getHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		alias := params["alias"]

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(alias))

			if b == nil {
				http.NotFound(w, r)
				return nil
			}

			c, _ := binary.Varint(b.Get([]byte("created")))
			u, _ := binary.Varint(b.Get([]byte("updated")))
			i, _ := binary.Varint(b.Get([]byte("id")))

			e := entry{string(alias),  i, c, u}

			js, err := json.Marshal(e)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

			return nil
		})
	}
}

func findHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func addHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e entry
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&e)

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte(e.Url))

			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}

			cb, ub, ib := make([]byte, 8), make([]byte, 8), make([]byte, 8)
			binary.PutVarint(cb, e.Created)
			binary.PutVarint(ub, e.Updated)
			binary.PutVarint(ib, e.Id)

			b.Put([]byte("created"), cb)
			b.Put([]byte("updated"), ub)
			b.Put([]byte("id"), ib)

			return nil
		})
	}
}
