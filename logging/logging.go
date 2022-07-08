package logging

/**
 * @Author: yc
 * @Description: 支持zap的logging
 * @File: logging
 * @Date: 2022/7/8 13:54
 */

import "go.uber.org/zap"

var logger *zap.Logger = nil


func SetLogger(l *zap.Logger) {
	logger = l
}


//zap的高性能接口,依赖zap的Field
func ZDebug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

func ZInfo(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}

}

func ZWarn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

func ZError(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

func ZFatal(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}

