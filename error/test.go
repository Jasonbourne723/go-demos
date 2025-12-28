package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidPath = errors.New("Invalid Path")
	ErrInternal    = errors.New("Internal Error")
)

type ConfigError struct {
	Err  error
	Path string
}

func (c *ConfigError) Error() string {
	if c.Err != nil {
		return fmt.Sprintf("config err,Path:%v,Cause of Error:%v", c.Path, c.Err)
	}
	return fmt.Sprintf("Invalid Path,Path: %v", c.Path)
}

func (c *ConfigError) Unwrap() error {
	return c.Err
}

type repo struct {
}

func NewRepo() *repo {
	return &repo{}
}

func (r *repo) ReadConfig(path string) (string, error) {
	fmt.Printf("recive path is %s\n", path)
	if len(path) == 0 {
		return "", &ConfigError{
			Path: path,
			Err:  ErrInvalidPath,
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return "", &ConfigError{
			Path: path,
			Err:  err,
		}
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", &ConfigError{
			Err:  err,
			Path: path,
		}
	}
	fmt.Printf("read data is %s\n", string(data))
	return string(data), nil
}

type service struct {
	repo *repo
}

func NewService(repo *repo) *service {
	return &service{
		repo,
	}
}

func (s *service) LoadConfig(path string) (string, error) {
	if data, err := s.repo.ReadConfig(path); err != nil {
		if errors.Is(err, ErrInvalidPath) {
			return "", ErrInvalidPath
		} else if errors.Is(err, os.ErrNotExist) {
			return "", os.ErrNotExist
		} else {
			return "", ErrInternal
		}
	} else {
		return data, nil
	}
}

type handler struct {
	srv *service
}

func NewHandler(srv *service) *handler {
	return &handler{
		srv,
	}
}

type Response struct {
	Data string
}

func (h *handler) LoadConfig(ctx *gin.Context) {
	path := ctx.Query("path")
	wd, _ := os.Getwd()
	fmt.Println(wd)
	fmt.Printf("read path is %v\n", path)
	data, err := h.srv.LoadConfig(path)
	fmt.Printf("http recive data is %s\n", data)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPath):
			ctx.AbortWithStatus(400)
		case errors.Is(err, os.ErrNotExist):
			ctx.AbortWithStatus(404)
		default:
			ctx.AbortWithStatus(500)
		}
	} else {
		ctx.JSON(200, &Response{
			Data: data,
		})
	}
}

func main() {

	//path := "./a.txt"

	engine := gin.Default()
	repo := NewRepo()
	srv := NewService(repo)

	//	data, err := srv.LoadConfig(path)
	//	if err != nil {
	//		fmt.Printf("err is %w\n", err)
	//	} else {
	//		fmt.Println(data)
	//	}

	hander := NewHandler(srv)

	engine.GET("/", hander.LoadConfig)
	go func() {

		err := engine.Run(":8080")
		if err != nil {
			fmt.Printf("httpserver start failed,err: %w\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("httpserver is closing")

}
