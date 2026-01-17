package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("please input fileName")
		return
	}
	fileName := flag.Args()[0]
	fmt.Printf("filename:%s\n", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if _, err := f.Seek(0, os.SEEK_END); err != nil {
		fmt.Println(err)
		return
	}
	buffers := make([]byte, 1024)
	for {
		stat, err := f.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}

		offset, err := f.Seek(0, os.SEEK_CUR)
		if err != nil {
			fmt.Println(err)
			return
		}
		if offset > stat.Size() {
			if _, err := f.Seek(0, os.SEEK_SET); err != nil {
				fmt.Println(err)
				return
			}
		}

		n, err := f.Read(buffers)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if n > 0 {
			fmt.Print(string(buffers[:n]))
		}
		<-time.After(time.Millisecond * 200)
	}
}
