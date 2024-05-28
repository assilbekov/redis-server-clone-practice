package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient1(t *testing.T) {
	c, err := NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SET val =>", fmt.Sprintf("Charlie_%d", 999))
	if err := c.Set(
		context.Background(),
		fmt.Sprintf("leader_%d", 999),
		"999 ",
	); err != nil {
		log.Fatal(err)
	}

	val, err := c.Get(context.Background(), fmt.Sprintf("leader_%d", 999))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET val =>", val)
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
