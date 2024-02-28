package main

import (
	torrentfile "bittorrent-client/internal/torrentFile"
	"crypto/rand"
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

	tf, _ := torrentfile.Open(logger, filePath)
	logger.Info("open")

	var peerID [20]byte
	rand.Read(peerID[:])

	fmt.Println(tf.RequestPeers(peerID, torrentfile.Port))

}
