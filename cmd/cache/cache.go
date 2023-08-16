package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/rshulabs/HgCache/internal/cache"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	cache.NewApp("hgcache").Run()
}
