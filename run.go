package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	carsAPI := exec.CommandContext(ctx, "node", "main.js")
	goBack := exec.CommandContext(ctx, "go", "run", ".")

	home, _ := os.Getwd()
	goBack.Dir = filepath.Join(home, "/go-backend")
	carsAPI.Dir = filepath.Join(home, "/cars-api")

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
		time.Sleep(time.Second * 2)
		if err := goBack.Run(); err != nil {
			log.Println("Go backend error:", err)
			cancel()
			os.Exit(1)
		}
	}()

	wg.Wait()
	log.Println()
}
