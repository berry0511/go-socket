package server

import (
    "github.com/natefinch/lumberjack"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger

func GetLogger() *zap.Logger {
    return zapLogger
}

func GetSugerLogger() *zap.SugaredLogger {
    return zapLogger.Sugar()
}

func NewLogger() {
    writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

    zapLogger = zap.New(core, zap.AddCaller())
}

func InitLogger(z *zap.Logger) {
    zapLogger = z
}

func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
    lumberJackLogger := &lumberjack.Logger{
        Filename:   "./Server.log",
        MaxSize:    100,
        MaxBackups: 10,
        MaxAge:     30,
        Compress:   true,
    }
    return zapcore.AddSync(lumberJackLogger)
}
