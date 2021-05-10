package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	_ "home_work/go/errgroup/dao"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	ctx     context.Context
	servers []*Server
	cancel  func()
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	helloSvr := NewServer(ServerName("hello"), Address("192.168.12.103"), Port(22334), AddHandle("/hello", Hello))
	byeSvr := NewServer(ServerName("bye"), Address("192.168.12.103"), Port(22335), AddHandle("/bye", Bye))
	return &App{
		ctx:     ctx,
		servers: []*Server{helloSvr, byeSvr},
		cancel:  cancel,
	}
}

func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done()
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func main() {
	fmt.Println("app")

	app := NewApp()
	err := app.Run()
	if err != nil {
		log.Fatal("app run err = ", err.Error())
	}
}
