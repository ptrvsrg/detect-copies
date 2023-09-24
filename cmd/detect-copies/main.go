package main

import (
	detectcopies "detect-copies/internal"
	cliparser "detect-copies/internal/cli-parser"
)

func main() {
	addr, port := cliparser.Parse()
	detectcopies.Run(addr, port)
}
