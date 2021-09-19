package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	httpSrv := http.NewServer(http.Address(":8080"))
	// http server
	g.Go(func() error {
		fmt.Println("http")
		go func() {
			<-ctx.Done()
			fmt.Println("http ctx done")
			httpSrv.Shutdown(context.TODO())
		}()
		return httpSrv.Start(context.Background())
	})

	// signal
	g.Go(func() error {
		// SIGTERM is POSIX specific
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
		sig := make(chan os.Signal, len(exitSignals))
		signal.Notify(sig, exitSignals...)
		for {
			fmt.Println("signal")
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case <-sig:
				// do something
				return nil
			}
		}
	})

	// inject error
	g.Go(func() error {
		fmt.Println("inject")
		time.Sleep(10 * time.Second)
		fmt.Println("inject finish")
		return errors.New("inject error")
	})

	// first error return
	err := g.Wait()
	fmt.Println(err)

}
