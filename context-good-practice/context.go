package main

import (
	"context"
	"fmt"
	"time"
)

/*
KEY IDEA:
Context is used for:
1) Cancellation
2) Deadlines
3) Request-scoped metadata (NOT business logic)
*/

// ----- context keys (best practice: unexported, unique type) -----

type requestIDKeyType struct{}
type userIDKeyType struct{}

var requestIDKey = requestIDKeyType{}
var userIDKey = userIDKeyType{}

// ----- entry point -----

func main() {
	// Root context (usually created by HTTP server / main)
	ctx := context.Background()

	// Add cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Add request-scoped data
	ctx = context.WithValue(ctx, requestIDKey, "req-123")
	ctx = context.WithValue(ctx, userIDKey, 42)

	// Simulate request handling
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println(">>> client canceled request")
		cancel()
	}()

	err := handleRequest(ctx)
	if err != nil {
		fmt.Println("request failed:", err)
	}
}

// ----- request handler -----

func handleRequest(ctx context.Context) error {
	log(ctx, "handling request")

	// Call deeper layers
	if err := authenticate(ctx); err != nil {
		return err
	}

	if err := fetchData(ctx); err != nil {
		return err
	}

	log(ctx, "request completed successfully")
	return nil
}

// ----- deeper service layer -----

func authenticate(ctx context.Context) error {
	log(ctx, "authenticating user")

	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func fetchData(ctx context.Context) error {
	log(ctx, "fetching data")

	select {
	case <-time.After(3 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ----- logging helper (GOOD use of context values) -----

func log(ctx context.Context, msg string) {
	reqID, _ := ctx.Value(requestIDKey).(string)
	userID, _ := ctx.Value(userIDKey).(int)

	fmt.Printf(
		"[requestID=%s userID=%d] %s\n",
		reqID,
		userID,
		msg,
	)
}
