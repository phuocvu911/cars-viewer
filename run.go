package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	carsAPI := exec.CommandContext(ctx, "node", "cars-api/main.js")
	goBack := exec.CommandContext(ctx, "go", "run", "./go-backend/main.go")

	carsAPI.Stdout = os.Stdout
	carsAPI.Stderr = os.Stderr
	goBack.Stdout = os.Stdout
	goBack.Stderr = os.Stderr

	go func() {
		defer wg.Done()
		fmt.Println("Starting Cars API...")
		if err := carsAPI.Run(); err != nil {
			log.Println("Cars API error:", err)
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Starting Go backend...")
		if err := goBack.Run(); err != nil {
			log.Println("Go backend error:", err)
			cancel()
			os.Exit(1)
		}
	}()

	wg.Wait()
}
