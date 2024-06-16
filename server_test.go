package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"redis-server-clone-practice/client"
	"sync"
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	listenAddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})

	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second / 2)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	testCases := map[string]string{
		"foo": "bar",
		"bar": "baz",
		"baz": "qux",
	}

	for k, v := range testCases {
		if err := rdb.Set(context.Background(), k, v, 0).Err(); err != nil {
			t.Fatal(err)
		}

		newVal, err := rdb.Get(context.Background(), k).Result()

		if err != nil {
			t.Fatal(err)
		}

		if newVal != v {
			t.Errorf("expected %s, got %s", v, newVal)
		}
	}

	fmt.Println("WE ARE HERE key", "value")
}

func TestFooBar(t *testing.T) {
	in := map[string]string{
		"server":  "redis",
		"version": "6.0 ",
	}
	out := respWriteMap(in)
	fmt.Println(string(out))
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
