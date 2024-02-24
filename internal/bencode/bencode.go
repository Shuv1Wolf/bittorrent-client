package bencode

import (
	"bittorrent-client/internal/lib/logger/sl"
	"log/slog"
	"os"
)

func Read(log *slog.Logger, filePath string) (string, error) {
	const op = "internal.bencode.read"

	log = log.With(slog.String("op", op))

	file, err := os.Open(filePath)

	if err != nil {
		log.Error("Ошибка при открытии файла", sl.Err(err))
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("Ошибка при получении информации о файле", sl.Err(err))
		return "", err
	}

	fileSize := fileInfo.Size()
	fileContent := make([]byte, fileSize)

	_, err = file.Read(fileContent)
	if err != nil {
		log.Error("Ошибка при чтении файла", sl.Err(err))
		return "", err
	}

	return string(fileContent), nil
}
