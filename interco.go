package main

import (
	"io"
	"net"
	"sync"

	"github.com/golang/snappy"
)

// interconnect c1 to c2 with optional snappy compression/decompression
// both c1 and c2 will be closed on completion
func interco(c1, c2 net.Conn, compress bool) {
	cl := make(chan struct{})
	once := sync.Once{}
	doClose := func() {
		once.Do(func() {
			close(cl)
		})
	}

	// closing both connections will trigger read failure and close the following goroutines
	defer c1.Close()
	defer c2.Close()

	// c1→c2
	go func() {
		defer doClose()
		if compress {
			// snappy compression
			w := snappy.NewBufferedWriter(c1)
			io.Copy(w, c2)
			w.Close()
		} else {
			io.Copy(c1, c2)
		}
	}()

	// c2→c1
	go func() {
		defer doClose()
		if compress {
			// snappy decompression
			io.Copy(c2, snappy.NewReader(c1))
		} else {
			io.Copy(c2, c1)
		}
	}()

	// wait for close signal
	<-cl
}
