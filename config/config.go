package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ENVIRONMENT    string
	MONGO_DB_NAME  string
	MONGO_USERNAME string
	MONGO_PASSWORD string
	MONGO_URL      string
	PORT           string
}

func GetConfiguration() Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf(getFileName(), &configuration)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}

	return configuration
}

func getFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "dev"
	}
	filename := []string{"/", "config.", env, ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
