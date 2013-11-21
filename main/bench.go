package main

import (
	"data"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	REQUESTS = 100000
	WORKERS  = 10
)

func client() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}
	enc, dec := json.NewEncoder(conn), json.NewDecoder(conn)
	for i := 0; i < REQUESTS; i++ {
		enc.Encode(data.M{Cmd: "set", Key: fmt.Sprintf("%d", i), Val: true})
		var m data.M
		dec.Decode(&m)
		if m.Val == false {
			log.Fatal("woops")
		}
	}
	conn.Close()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(WORKERS)
	for i := 0; i < WORKERS; i++ {
		go func() {
			client()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(WORKERS * REQUESTS)
}
