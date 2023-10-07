package engine

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/RobDoan/go-docker-template/pkg/services"
	"github.com/RobDoan/go-docker-template/pkg/utils"
	"github.com/RobDoan/go-docker-template/pkg/vault"
	"github.com/judedaryl/go-arrayutils"
	"github.com/pkg/errors"
)

const recursionMaxNums = 100

type BuildOptions struct {
	LocalServices       []string
	ServiceNeedLocalDBs []string
	Environment         string
	VaultUrl            string
	Channel             string
	CountryCode         string
	WorkDir             string
	connectionLinks     map[string]string
}

type rendereredValues struct {
	Name        string
	Environment map[string]interface{}
	DebugPorts  map[string]int
	Command     interface{}
}

func NewBuildOptions(services []string,
	environment string,
	vault string,
	channel string,
	countryCode string,
	workDir string) *BuildOptions {

	serviceNeedlocaDb := []string{"gibson"}
	path, err := utils.GetAbsDirectory(workDir)
	if err != nil {
		utils.PrintError("Unable to get absolute path")
		return nil
	}

	workingDirectory := filepath.Join(path, "docker-services")

	return &BuildOptions{
		LocalServices:       services,
		Environment:         environment,
		VaultUrl:            vault,
		Channel:             channel,
		CountryCode:         countryCode,
		WorkDir:             workingDirectory,
		ServiceNeedLocalDBs: serviceNeedlocaDb,
	}
}

func (b *BuildOptions) Build() error {
	serviceName := "gibson"
	serviceValues, err := services.ReadServiceValue(b.WorkDir, serviceName)
	if err != nil {
		utils.PrintError("Unable to read service configuartion values. Please check values.yaml file")
		return err
	}

	environmentVariables, err := b.getEnvironmentVarialbes(serviceName, serviceValues)
	if err != nil {
		utils.PrintError("Unable to get environment variables")
		return err
	}

	utils.PrintInfo("Copying database from remote to local")
	if b.isNeedALocalDb(serviceName) {
		utils.PrintInfo("Need to ask devops for permission to access and dump database to local")
		// target := getLocalImportDbConfig(serviceName, b.Channel)
		// source, err := getDBConfigValue(environmentVariables, b.Channel)
		// if err != nil {
		// 	utils.PrintError("Unable to get database config value")
		// 	return err
		// }
		// if err := loader.CopyDatabase(&target, &source); err != nil {
		// 	utils.PrintError("Unable to copy database from remote to local")
		// 	return err
		// }
		// COPY DATABASE FROM REMOTE TO LOCAL
		// update environment variables
	}

	values := rendereredValues{
		Name:        serviceName,
		Environment: environmentVariables,
		DebugPorts:  serviceValues.DebugPorts,
		Command:     serviceValues.Command,
	}

	if err := b.Render(serviceName, values); err != nil {
		utils.PrintError("Unable to render service")
		return err
	}
	return nil
}

func (b *BuildOptions) isNeedALocalDb(serviceName string) bool {
	return arrayutils.Some(b.ServiceNeedLocalDBs, func(val string) bool { return val == serviceName })
}

func (b *BuildOptions) readConnections(yamlFile string) (map[string]string, error) {
	if b.connectionLinks != nil {
		return b.connectionLinks, nil
	}

	connections := make(map[string]map[string]string)
	if err := utils.ReadYamlFile(yamlFile, &connections); err != nil {
		return nil, err
	}
	connectionLinks := map[string]string{}
	for serviceName, v := range connections {
		connectionLinks[serviceName] = v[b.Environment]
	}
	b.connectionLinks = connectionLinks
	return connectionLinks, nil
}

func (b *BuildOptions) connectUrl(serviceName string) (string, error) {
	serviceConnectionFile := filepath.Join(b.WorkDir, "service-connections.yaml")
	connections, err := b.readConnections(serviceConnectionFile)
	if err != nil {
		utils.PrintError("Unable to read service-connections.yaml file")
		return "", err
	}
	if _, ok := connections[serviceName]; !ok {
		utils.PrintError(fmt.Sprintf("Unable to find connection for service %s", serviceName))
		return "", err
	}
	return connections[serviceName], nil
}

func (b *BuildOptions) readVaults(service string) (map[string]interface{}, error) {
	vault := vault.NewVault(b.VaultUrl)
	secret, err := vault.GetVaults(b.Environment, service)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return secret.Data, nil
}
func removeKeysMatch(re regexp.Regexp, m map[string]interface{}) map[string]interface{} {
	for k := range m {
		if re.MatchString(k) {
			delete(m, k)
		}
	}
	return m
}

func ingoreServerEnvVars(m map[string]interface{}) map[string]interface{} {
	for _, envName := range IgnoreUnusedEnvVars {
		delete(m, envName)
	}
	return m
}

func (b *BuildOptions) getEnvironmentVarialbes(service string, serviceValues services.ServiceValues) (map[string]interface{}, error) {
	environmentData, err := b.readVaults(service)
	if err != nil {
		utils.PrintError("Unable to read vaults")
		fmt.Println(err)
		return nil, err
	}

	for serviceName, env := range serviceValues.LinkEnvironmentName {
		isLocalService := arrayutils.Some(b.LocalServices, func(val string) bool { return val == serviceName })
		if isLocalService {
			utils.PrintInfo(fmt.Sprintf("[%s] will link with local container", serviceName))
			delete(environmentData, env)
			continue
		}

		if serviceName == "mysql" {
			utils.PrintInfo("Skipping setting value for mysql service")
			continue
		}

		connectionUrl, err := b.connectUrl(serviceName)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Unable to get the connect url for [%s] service", serviceName))
			return nil, err
		}
		environmentData[env] = connectionUrl
	}
	data := removeKeysMatch(*regexp.MustCompile(`^REDIS_*`), environmentData)
	data = ingoreServerEnvVars(data)
	return data, nil
}

// This function is copied from helm source code
func initFunMap(tpl *template.Template) {
	funcMap := funcMap()
	includedNames := make(map[string]int)

	funcMap["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", errors.Wrapf(fmt.Errorf("unable to execute template"), "rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := tpl.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}
	funcMap["getPort"] = func(portsConfig map[string]int, name string) int {
		return portsConfig[name]
	}

	tpl.Funcs(funcMap)
}

func (b *BuildOptions) Render(serviceName string, values rendereredValues) error {
	helperTemplatePath := filepath.Join(b.WorkDir, serviceName, services.HelperfileName)
	serviceTemplatePath := filepath.Join(b.WorkDir, serviceName, services.ServiceTemplateFileName)
	t := template.New("docker-compose.yml")
	initFunMap(t)
	tmp, err := t.ParseFiles(helperTemplatePath, serviceTemplatePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Unable to parse template of [%s] service", serviceName))
		fmt.Println(err)
		return err
	}
	var str strings.Builder
	if err := tmp.ExecuteTemplate(&str, services.ServiceTemplateFileName, values); err != nil {
		log.Fatal("Unable to execute template: ", err)
		return err
	}
	final := strings.ReplaceAll(str.String(), "<no value>", "")
	utils.WriteFile(filepath.Join(b.WorkDir, serviceName, "docker-compose.override.yaml"), []byte(final))
	return nil
}
