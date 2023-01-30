package services

import (
	"fmt"
	"jwt/pkg/helpers/utils"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(utils.ResolveProjectPath("jwt") + `/deployments/.env`)
	if err != nil {
		fmt.Println(".env file isn't found")
	}
}
