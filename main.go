package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Response struct {
	value int
	err   error
}

func fetchSlowThirdPartyStuff() (int, error) {
	time.Sleep(time.Millisecond * 1000)
	return 201, nil
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	value := ctx.Value("foo")
	fmt.Println("value in ctx:-", value)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	response := make(chan Response)

	go func() {
		val, err := fetchSlowThirdPartyStuff()
		response <- Response{val, err}
	}()

	// for {
	select {
	case <-ctx.Done(): // ?We write to Done channel when the timeout is triggered
		return -1, fmt.Errorf("Api took to long to respond")
	case r := <-response:
		return r.value, r.err
	}
	// }
}

func main() {
	start := time.Now()
	// ctx := context.Background()
	ctx := context.WithValue(context.Background(), "foo", "bar")
	userID := 10
	val, err := fetchUserData(ctx, userID)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:-", val)
	fmt.Println(time.Since(start))
}
