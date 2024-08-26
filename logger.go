package time

import (
	"github.com/ameise84/logger"
)

type Logger interface {
	Error(any)
}

var _gLogger Logger

func init() {
	_gLogger = logger.DefaultLogger()
}

func SetLogger(log logger.Logger) {
	_gLogger = log
}
