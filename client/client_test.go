package client

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Set(context.Background(), "foo", "bar", 0).Err(); err != nil {
		t.Fatal(err)
	}

	if _, err := rdb.Get(context.Background(), "foo").Result(); err != nil {
		t.Fatal(err)
	}

	fmt.Println("WE ARE HERE key", "value")

	/*val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)*/
}

func TestNewClient1(t *testing.T) {
	c, err := NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	defer func(c *Client) {
		err := c.Close()
		if err != nil {

		}
	}(c)

	if err := c.Set(
		context.Background(),
		"leader",
		1,
	); err != nil {
		log.Fatal(err)
	}

	val, err := c.Get(context.Background(), "leader")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET val =>", val)

	time.Sleep(time.Second * 2)
	// select {} // we are blocking here to keep the server running
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
