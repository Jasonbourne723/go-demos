package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var path *string = flag.String("path", "", "")

func main() {

	flag.Parse()

	if txt, err := ReadConfig(*path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("custom err: %s\n", err.Error())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("read successed,data: %s", txt)
	}

}

type ConfigError struct {
	Err  error
	Path string
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Invail Path: %s,Error cause %v", e.Path, e.Err)
	}
	return fmt.Sprintf("Invail Path: %s", e.Path)
}

func (e *ConfigError) Unwrap1() error {
	return e.Err
}

func ReadConfig(path string) (string, error) {

	if len(path) == 0 {
		return "", &ConfigError{
			Path: path,
			Err:  errors.New("Empty Path"),
		}
	}

	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", &ConfigError{
				Err:  err,
				Path: path,
			}
		} else {
			return "", &ConfigError{
				Err:  err,
				Path: path,
			}
		}
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", &ConfigError{
			Path: path,
			Err:  fmt.Errorf("read file failed,err: %w", err),
		}
	} else {
		return string(b), nil
	}
}
