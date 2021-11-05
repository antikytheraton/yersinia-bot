package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/antikytheraton/yersinia-bot/bot"
	"github.com/antikytheraton/yersinia-bot/downloader"
	log "github.com/sirupsen/logrus"
)

const usage = `Usage %s command [args]

Commands:
    yersinia-bot	- run telegram chatbot 
`

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}

	flag.Parse()

	var err error

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	dl := downloader.NewYtDownloader()

	switch flag.Arg(0) {
	case "yersinia-bot":
		err = bot.Run(dl, flag.Args()[1:])
	default:
		flag.Usage()
		return 1
	}

	if err != nil {
		log.Error(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
