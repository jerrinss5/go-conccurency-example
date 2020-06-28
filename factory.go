package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Inserter1 ... is used to place items to the channel aka belt
func Inserter1(belt chan int) {
	i := 0
	for i < 5 {
		time.Sleep(1 * time.Second)
		belt <- i
		i++
	}
	fmt.Println("[Inserter 1] I have finished placing items on the belt. Closing the belt")
	close(belt)
	return
}

// Inserter2 ... is used to pick items from the channel aka belt
func Inserter2(ctx context.Context, belt chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case i, ok := <-belt:
			if ok {
				fmt.Printf("[Inserter 2] I got %d\n", i)
			} else {
				fmt.Println("[Inserter 2] Belt has been closed. No more items to pick!")
				return
			}
		case <-ctx.Done():
			fmt.Printf("[Inserter 2] Cancelling work because : %q\n", ctx.Err())
			return
		}
	}
}

// KillSwitch ... to capture control + C from the user can call cancel
func KillSwitch(cancel context.CancelFunc) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("[Kill Switch] Ctrl-C pressed. Cancelling everything")
	cancel()
}

func main() {
	belt := make(chan int)
	var waitGroup sync.WaitGroup

	go Inserter1(belt)
	backgroundCtx := context.Background()
	ctx, cancel := context.WithTimeout(backgroundCtx, 3*time.Second)
	go KillSwitch(cancel)
	for i := 0; i < 3; i++ {
		waitGroup.Add(1)
		go Inserter2(ctx, belt, &waitGroup)
	}
	waitGroup.Wait()
	fmt.Println("[Main] All Done! Exiting!")
}
