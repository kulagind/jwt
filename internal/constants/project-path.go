package constants

import (
	"jwt/pkg/helpers/utils"
	"os"
)

func ProjectPath() string {
	root := "/"
	if os.Getenv("APP_MODE") != "prod" {
		root = utils.ResolveProjectPath("jwt")
	}
	return root
}
