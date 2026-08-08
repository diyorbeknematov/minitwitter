package logger

import (
	"io"
	"log/slog"
	"os"
)

func SetupLog() *slog.Logger {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	return slog.New(
		slog.NewJSONHandler(multiWriter, nil),
	)
}
