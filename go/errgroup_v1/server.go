package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

// option
type ServerOption func(*Server)

func ServerName(name string) ServerOption {
	return func(s *Server) {
		s.name = name
	}
}

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func Port(p int) ServerOption {
	return func(s *Server) {
		s.port = strconv.Itoa(p)
	}
}

func AddHandle(name string, h ServerHandleFunc) ServerOption {
	return func(s *Server) {
		handle := &ServerHandle{
			name: name,
			fn:   h,
		}
		s.handles = append(s.handles, handle)
	}
}

// handles
type ServerHandle struct {
	name string
	fn   ServerHandleFunc
}

type ServerHandleFunc func(http.ResponseWriter, *http.Request)

func Hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello")
}

func Bye(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "bye")
}

// server
type Server struct {
	name    string
	address string
	port    string
	handles []*ServerHandle
}

func NewServer(opt ...ServerOption) *Server {
	srv := &Server{
		address: "",
		port:    "",
	}
	for _, o := range opt {
		o(srv)
	}
	return srv
}

func (s *Server) Start() error {
	for _, handle := range s.handles {
		http.HandleFunc(handle.name, handle.fn)
	}
	err := http.ListenAndServe(s.address+":"+s.port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	log.Println("")
	return nil
}
