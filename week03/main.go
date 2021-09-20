package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	// http server
	eg.Go(func() error {
		server := http.Server{Addr: ":8080"}

		go func() {
			// 启动后监听信号
			select {
			case <-ctx.Done():
				fmt.Println("http context done")
				err := server.Shutdown(ctx)
				if err != nil {
					log.Fatal("stop http server fail")
				}
			}
		}()

		fmt.Println("http server start")
		return server.ListenAndServe()
	})

	// signal
	eg.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL}
		sigChan := make(chan os.Signal, len(exitSignals))
		signal.Notify(sigChan, exitSignals...)
		fmt.Println("waiting signal...")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("signal context done")
				return ctx.Err()
			case sig := <-sigChan:
				fmt.Printf("receive quit command -> %s\n", sig)
				//收到信号后，将context标记为done
				ctx.Done()
				return errors.Errorf("receive quit command -> %s", sig)
			}
		}
	})

	err := eg.Wait()
	fmt.Printf("main stop -> %s", err)

}
