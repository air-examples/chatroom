package gas

import "github.com/aofei/air"

var (
	Auth      = AuthHandler()
	PreRender = PreRenderHandler()
	PreLogger = PreLoggerHandler()
)

func InitGas() {
	air.Pregases = append(air.Pregases, PreLogger)
}
