package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"sync"
)

func main() {
	var fileName string
	var chunk int
	flag.StringVar(&fileName, "file", "", "-file=name")
	flag.IntVar(&chunk, "chunk", -1, "chunk size in MB, -chunk=5")
	flag.Parse()

	if fileName == "" {
		log.Fatal("-file flag is required")
	}

	if chunk <= 0 {
		log.Fatal("chunk size must be greater than zero")
	}

	err := run(fileName, chunk)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("file processing done, bye, bye")

}

func run(fileName string, chunk int) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("file %s size %d", file.Name(), info.Size())

	var chunkNum = 1
	var chunkSize = chunk << 20

	pool := sync.Pool{}
	pool.New = func() interface{} {
		return make([]byte, chunkSize)
	}

	reader := bufio.NewReader(file)

	g := errgroup.Group{}

	for {
		// A Pool is safe for use by multiple goroutines simultaneously. https://pkg.go.dev/sync#Pool
		buff := pool.Get().([]byte)
		cnt, err := reader.Read(buff)

		if err == io.EOF {
			log.Println("file processing done")
			break
		}

		if err != nil {
			log.Println("reader error: ", err)
			break
		}

		log.Printf("chunk %d  read %d bytes", chunkNum, cnt)

		g.Go(func() error {
			defer pool.Put(buff)

			chunkFile, err := os.Create(fmt.Sprintf("chunk_%d_%s", chunkNum, info.Name()))
			if err != nil {
				log.Println("error creating chunk file", err)
				return err
			}
			defer chunkFile.Close()

			writer := bufio.NewWriter(chunkFile)
			defer writer.Flush()

			cnt, err = writer.Write(buff[:cnt])
			if err != nil {
				log.Println("writer error: ", err)
				return err
			}

			log.Printf("chunk %d write %d bytes", chunkNum, cnt)
			return nil
		})

		if err := g.Wait(); err != nil {
			log.Println("error processing chunk: ", err)
			log.Fatal(err)
		}

		if chunkNum > int(info.Size())/chunkSize {
			break
		}

		chunkNum++
	}
	return nil
}
