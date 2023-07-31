package defaultLogger

import (
	"os"

	"github.com/danilluk1/social-network/libs/go/pkg/constants"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/models"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/zap"
)

var Logger logger.Logger

func SetupDefaultLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		Logger = zap.NewZapLogger(
			&conf.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.Dev,
		)
		break
	// case "Logrus":
	// 	Logger = logrous.NewLogrusLogger(
	// 		&config.LogOptions{LogType: models.Logrus, CallerEnabled: false},
	// 		constants.Dev,
	// 	)
	// 	break
	default:
	}
}
