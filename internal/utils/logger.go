package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	// Создаем новый логгер
	Log = logrus.New()

	// Настраиваем формат вывода
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Настраиваем уровень логирования
	Log.SetLevel(logrus.DebugLevel)

	// Выводим логи в stdout (консоль)
	Log.SetOutput(os.Stdout)
}
