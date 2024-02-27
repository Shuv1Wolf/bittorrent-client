package bencode

import (
	"bittorrent-client/internal/lib/logger/sl"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func Open(log *slog.Logger, filePath string) (TorrentFile, error) {
	const op = "internal.bencode.Open"
	log = log.With(slog.String("op", op))

	file, err := os.Open(filePath)

	if err != nil {
		log.Error("error when opening file", sl.Err(err))
		return TorrentFile{}, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		log.Error("error when Unmarshal", sl.Err(err))
		return TorrentFile{}, fmt.Errorf("%s: %w", op, err)
	}
	return bto.toTorrentFile()
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	const op = "internal.bencode.hash"

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, fmt.Errorf("%s: %w", op, err)
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	const op = "internal.bencode.splitPieceHashes"

	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		return nil, fmt.Errorf("%s: received malformed pieces of length %d", op, len(buf))
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	const op = "internal.bencode.toTorrentFile"

	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, fmt.Errorf("%s: %w", op, err)
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, fmt.Errorf("%s: %w", op, err)
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}
