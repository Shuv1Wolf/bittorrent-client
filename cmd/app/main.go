package main

import (
	"bittorrent-client/internal/bencode"
	"flag"
	"fmt"
	"log/slog"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file-path", "", "path to bittorrent file")
	flag.Parse()

	if filePath == "" {
		panic("file-path is required")
	}

	logger := slog.New(slog.Default().Handler())

	logger.Info("start app")

	fmt.Println(bencode.Read(logger, filePath))
	bencode.Read(logger, filePath)
}
