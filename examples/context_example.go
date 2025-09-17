package main

import (
	"context"
	"fmt"
	"time"
)

// Example showing different context patterns

// 1. Context with timeout
func databaseOperationWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always call cancel to free resources

	// Simulate a database operation
	select {
	case <-time.After(3 * time.Second): // This takes 3 seconds (longer than timeout)
		fmt.Println("Database operation completed")
	case <-ctx.Done():
		fmt.Println("Database operation timed out:", ctx.Err())
	}
}

// 2. Context with cancellation
func longRunningTask(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled:", ctx.Err())
			return ctx.Err()
		default:
		}

		// Do some work
		fmt.Printf("Working... step %d\n", i+1)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Task completed successfully")
	return nil
}

// 3. Context with values (like user ID)
func processWithUserContext() {
	ctx := context.Background()

	// Add user ID to context (like your middleware does)
	ctx = context.WithValue(ctx, "userID", int64(123))
	ctx = context.WithValue(ctx, "username", "john_doe")

	// Pass context down the call chain
	handleUserRequest(ctx)
}

func handleUserRequest(ctx context.Context) {
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		fmt.Println("No user ID in context")
		return
	}

	username, ok := ctx.Value("username").(string)
	if !ok {
		fmt.Println("No username in context")
		return
	}

	fmt.Printf("Processing request for user %s (ID: %d)\n", username, userID)
}

func main() {
	fmt.Println("=== Context with Timeout ===")
	databaseOperationWithTimeout()

	fmt.Println("\n=== Context with Cancellation ===")
	ctx, cancel := context.WithCancel(context.Background())

	// Start a long-running task
	go longRunningTask(ctx)

	// Cancel it after 2 seconds
	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second) // Wait to see the cancellation

	fmt.Println("\n=== Context with Values ===")
	processWithUserContext()
}
