package main

import (
	"runtime"

	"github.com/ccxt-go/ccxt-go/cmd"
)

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute()
}
