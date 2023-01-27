package utils

import (
	"os"
	"regexp"
)

func ResolveProjectPath(projectDirName string) string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	return string(rootPath)
}
