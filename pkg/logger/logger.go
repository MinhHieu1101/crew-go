package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	l, _ := zap.NewDevelopment()
	Log = l
}
