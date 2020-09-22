package server

import (
    "github.com/natefinch/lumberjack"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "sync"
)

var sugarLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
    return sugarLogger
}

var once sync.Once

func GetSingletonObj() *zap.SugaredLogger {
    once.Do(func() {
        sugarLogger = &zap.SugaredLogger{}
    })
    return sugarLogger
}

func InitLogger() {
    writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

    logger := zap.New(core, zap.AddCaller())
    sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
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
