package main

import (
	"fmt"

	detectcopies "detect-copies/internal"
	cliparser "detect-copies/internal/cli-parser"
)

const title = "                                                               \n" +
	"    ____       __            __                      _         \n" +
	"   / __ \\___  / /____  _____/ /_   _________  ____  (_)__  _____\n" +
	"  / / / / _ \\/ __/ _ \\/ ___/ __/  / ___/ __ \\/ __ \\/ / _ \\/ ___/\n" +
	" / /_/ /  __/ /_/  __/ /__/ /_   / /__/ /_/ / /_/ / /  __(__  )\n" +
	"/_____/\\___/\\__/\\___/\\___/\\__/   \\___/\\____/ .___/_/\\___/____/ \n" +
	"                                          /_/                  \n" +
	"Detect-Copies: 1.0.0                                           \n"

func main() {
	fmt.Println(title)
	addr, port := cliparser.Parse()
	detectcopies.Start(addr, port)
}
