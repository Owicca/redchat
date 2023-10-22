package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	addr  = "localhost:6379"
	user  = "root"
	pass  = "root"
	user1 = "user"
	pass1 = "user"
)

var (
	stream   = "stream1"
	isClient = false
	isServer = true
	wg       sync.WaitGroup
)

func main() {
	if len(os.Args) > 1 {
		stream = os.Args[1]
	}

	readCLi()

	if isClient {
		wg.Add(1)
		go client(wg)
	} else if isServer {
		wg.Add(1)
		go server(wg)
	} else {
		wg.Add(1)
		go server(wg)

		wg.Add(1)
		go client(wg)
	}

	wg.Wait()

	log.Println("[INFO]", "END")
}

func readCLi() {
	flag.BoolVar(&isClient, "client", false, "Is client (bool)")
	flag.BoolVar(&isServer, "server", true, "Is server (bool)")
	flag.StringVar(&stream, "stream", "stream1", "The name of the stream")

	if isClient {
		isServer = false
	} else if isServer {
		isClient = false
	}

	flag.Parse()
}

func server(wg sync.WaitGroup) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pass,
	})

	ctx := context.Background()

	for {
		val := uuid.NewString()
		id, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: "stream1",
			ID:     "*",
			Values: map[string]interface{}{"key": val},
		}).Result()
		if err != nil {
			wg.Done()
			log.Fatal("[ERROR]", err)
		}

		delta := rand.Intn(1000)
		log.Println("[INFO]", "ADD", stream, user, delta, id, val)

		time.Sleep(time.Duration(delta) * time.Millisecond)
	}

	log.Println("[INFO]", "server end")
	wg.Done()
}

func client(wg sync.WaitGroup) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user1,
		Password: pass1,
	})

	ctx := context.Background()

	for {
		streamList, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, "$"},
		}).Result()
		if err != nil {
			wg.Done()
			log.Fatal("[ERROR]", err)
		}

		for _, strm := range streamList {
			for _, msg := range strm.Messages {
				for k, v := range msg.Values {
					log.Println("[INFO]", "READ", stream, user1, k, v)
				}
			}
		}
	}

	log.Println("[INFO]", "client end", client, ctx)
	wg.Done()
}
