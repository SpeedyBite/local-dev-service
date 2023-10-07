package engine

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/RobDoan/go-docker-template/pkg/loader"
	"github.com/RobDoan/go-docker-template/pkg/utils"
)

var DB_READ_SUFFIX = []string{
	"DB_READ_HOST",
	"DB_READ_USERNAME",
	"DB_READ_PASSWORD",
}
var DB_WRITE_SUFFIX = []string{
	"DB_WRITE_HOST",
	"DB_WRITE_USERNAME",
	"DB_WRITE_PASSWORD",
}

const DB_NAME_SUFFIX = "DB_DATABASE"
const localImportDbName = "tmp_dev"
const localDbPassword = "payfare-dev"
const localDbUser = "payfare-dev"

func getLocalImportDbConfig(service, channel string) loader.DatabaseConfig {
	databaseName := fmt.Sprintf("%s_%s_%s", channel, service, localImportDbName)
	return loader.DatabaseConfig{
		Host:     "localhost",
		User:     localDbUser,
		Password: localDbPassword,
		Name:     databaseName,
	}
}

func getDBConfigValue(environmentData map[string]interface{}, channel string) (loader.DatabaseConfig, error) {
	configVal := map[string]string{}
	dbConfig := loader.DatabaseConfig{}

	databaseEnvVarName := fmt.Sprintf("%s_%s", strings.ToUpper(channel), DB_NAME_SUFFIX)
	if value, ok := environmentData[databaseEnvVarName]; ok {
		configVal[DB_NAME_SUFFIX] = value.(string)
	}

	for _, suffix := range DB_WRITE_SUFFIX {
		envVarName := fmt.Sprintf("%s_%s", strings.ToUpper(channel), suffix)
		if value, ok := environmentData[envVarName]; ok {
			configVal[suffix] = value.(string)
		}
	}

	if len(configVal) == 0 {
		utils.PrintError("No database config found")
		return dbConfig, errors.New("No database config found")
	}
	dbConfig.Host = configVal["DB_WRITE_HOST"]
	dbConfig.User = configVal["DB_WRITE_USERNAME"]
	dbConfig.Password = configVal["DB_WRITE_PASSWORD"]
	dbConfig.Name = configVal[DB_NAME_SUFFIX]

	return dbConfig, nil
}

func updateEnvironmentVariables(environmentData *map[string]interface{}, channel string, dbConfig loader.DatabaseConfig) {
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), DB_NAME_SUFFIX)] = dbConfig.Name
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_READ_HOST")] = dbConfig.Host
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_WRITE_HOST")] = dbConfig.Host
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_READ_USERNAME")] = dbConfig.User
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_WRITE_USERNAME")] = dbConfig.User
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_READ_PASSWORD")] = dbConfig.Password
	(*environmentData)[fmt.Sprintf("%s_%s", strings.ToUpper(channel), "DB_WRITE_PASSWORD")] = dbConfig.Password

}
