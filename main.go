package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// Create an errgroup with a context.
	// If any goroutine returns an error, 'ctx' will be cancelled.
	g, ctx := errgroup.WithContext(context.Background())

	// Task 1: Successful task
	g.Go(func() error {
		// Simulate work
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Task 1 completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Task 1 cancelled")
			return ctx.Err()
		}
	})

	// Task 2: Failing task
	g.Go(func() error {
		// Simulate a quick failure
		time.Sleep(100 * time.Millisecond)
		return fmt.Errorf("Task 2 failed critically!")
	})

	// Task 3: Another task (will be cancelled because Task 2 fails)
	g.Go(func() error {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Task 3 completed")
			return nil
		case <-ctx.Done():
			// This happens because Task 2 failed
			fmt.Println("Task 3 cancelled because sibling failed")
			return ctx.Err()
		}
	})

	// Wait for all to finish.
	// Returns the first error encountered (from Task 2).
	if err := g.Wait(); err != nil {
		fmt.Printf("Outcome: The group operation failed: %v\n", err)
	} else {
		fmt.Println("Outcome: All tasks succeeded!")
	}
}
