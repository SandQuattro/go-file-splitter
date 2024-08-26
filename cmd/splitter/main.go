package main

import (
	"bufio"
	"flag"
	"fmt"
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
	log.Println("file processing done")

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

	for {
		chunkFile, err := os.Create(fmt.Sprintf("chunk_%d_%s", chunkNum, info.Name()))
		if err != nil {
			log.Println("error creating chunk file", err)
			return err
		}
		writer := bufio.NewWriter(chunkFile)

		buff := pool.Get().([]byte)
		cnt, err := reader.Read(buff)

		if err == io.EOF {
			writer.Flush()
			chunkFile.Close()
			log.Println("file processing done")
			break
		}

		if err != nil {
			log.Println("reader error: ", err)
			break
		}

		log.Printf("chunk %d read %d bytes", chunkNum, cnt)

		cnt, err = writer.Write(buff[:cnt])
		if err != nil {
			log.Println("writer error: ", err)
			return err
		}

		log.Printf("chunk %d write %d bytes", chunkNum, cnt)

		writer.Flush()
		chunkFile.Close()

		pool.Put(buff)

		if chunkNum > int(info.Size())/chunkSize {
			break
		}

		chunkNum++
	}
	return nil
}
