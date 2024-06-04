package main

import (
	"context"
	"fmt"
	"log"
	"redis-server-clone-practice/client"
	"sync"
	"testing"
	"time"
)

func TestFooBar(t *testing.T) {
	in := map[string]string{
		"foo":  "bar",
		"baz":  "qux",
		"quux": "corge",
	}
	out := respWriteMap(in)
	fmt.Println(out)
}

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := 0; i < nClients; i++ {
		i := i
		go func() {
			defer wg.Done()
			c, err := client.NewClient("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}

			defer func(c *client.Client) {
				err := c.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(c)

			key := fmt.Sprintf("client_key_%d", i)
			value := fmt.Sprintf("client_val_%d", i)
			if err := c.Set(
				context.Background(),
				key,
				value,
			); err != nil {
				log.Fatal(err)
			}

			val, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("GET val => %s, FROM client => %s\n", val, key)
		}()
	}

	wg.Wait()

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Errorf("expected 0 peers, got %d", len(server.peers))
	}
}
