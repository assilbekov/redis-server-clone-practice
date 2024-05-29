package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := 0; i < nClients; i++ {
		i := i
		go func() {
			c, err := NewClient("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}

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

			wg.Done()
		}()
	}

	wg.Wait()
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("SET val =>", fmt.Sprintf("Charlie_%d", i))
		if err := c.Set(
			context.Background(),
			fmt.Sprintf("leader_%d", i),
			fmt.Sprintf("Charlie_%d", i),
		); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.Background(), fmt.Sprintf("leader_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET val =>", val)
	}

	time.Sleep(time.Second * 2)
	// select {} // we are blocking here to keep the server running
}
