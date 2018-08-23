package setting

import (
	"gopkg.in/ini.v1"
	"time"
	"path"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"net/url"

	"regexp"

	"github.com/emretiryaki/merkut/pkg/util"
	"github.com/emretiryaki/merkut/pkg/log"

	"github.com/go-macaron/session"
)

type Scheme string

const (
	HTTP              Scheme = "http"
	HTTPS             Scheme = "https"
	SOCKET            Scheme = "socket"
	DEFAULT_HTTP_ADDR string = "0.0.0.0"
)

const (
	DEV  string = "development"
	PROD string = "production"
	TEST string = "test"
)
var (
	// App settings.
	Env          = DEV
	AppUrl       string
	AppSubUrl    string
	ApplicationName string
	InstanceName string


	BuildStamp      int64

	// Paths
	LogsPath       string
	HomePath       string
	DataPath       string
	PluginsPath    string
	CustomInitPath = "conf/custom.ini"

	// Global setting objects.
	Raw          *ini.File
	ConfRootPath string

	// Http server options
	Protocol           Scheme
	Domain             string
	HttpAddr, HttpPort string
	SshPort            int
	CertFile, KeyFile  string
	SocketPath         string
	RouterLogging      bool
	DataProxyLogging   bool
	StaticRootPath     string
	EnableGzip         bool
	EnforceDomain      bool

	DashboardVersionsToKeep int

	// for logging purposes
	configFiles                  []string
	appliedCommandLineProperties []string
	appliedEnvOverrides          []string

	// Log settings.
	LogModes   []string
	LogConfigs []util.DynMap

	AdminUser string
	VerifyEmailEnabled      bool

	// Security settings.
	SecretKey                        string
	LogInRememberDays                int
	CookieUserName                   string
	CookieRememberName               string
	DisableGravatar                  bool
	EmailCodeValidMinutes            int
	DataProxyWhiteList               map[string]bool
	DisableBruteForceLoginProtection bool
	SignoutRedirectUrl      string

	// Session settings.
	SessionOptions         session.Options
	SessionConnMaxLifetime int64
)

type Cfg struct {

	Raw *ini.File
	// Paths
	ProvisioningPath string

	// SMTP email settings
	Smtp SmtpSettings

	// Rendering
	ImagesDir                        string
	PhantomDir                       string
	RendererUrl                      string
	DisableBruteForceLoginProtection bool

	TempDataLifetime time.Duration
}

func NewCfg() *Cfg {
	return &Cfg{
		Raw: ini.Empty(),
	}
}

type CommandLineArgs struct {
	Config   string
	HomePath string
	Args     []string
}


func setHomePath(args *CommandLineArgs) {
	if args.HomePath != "" {
		HomePath = args.HomePath
		return
	}

	HomePath, _ = filepath.Abs(".")
	// check if homepath is correct
	if pathExists(filepath.Join(HomePath, "conf/defaults.ini")) {
		return
	}

	// try down one path
	if pathExists(filepath.Join(HomePath, "../conf/defaults.ini")) {
		HomePath = filepath.Join(HomePath, "../")
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (cfg *Cfg) Load(args *CommandLineArgs) error {
	setHomePath(args)

	iniFile, err := loadConfiguration(args)
	if err != nil {
		return err
	}

	cfg.Raw = iniFile

	// Temporary keep global, to make refactor in steps
	Raw = cfg.Raw

	ApplicationName = "Merkur"


	Env = iniFile.Section("").Key("app_mode").MustString("development")
	InstanceName = iniFile.Section("").Key("instance_name").MustString("unknown_instance_name")
	PluginsPath = makeAbsolute(iniFile.Section("paths").Key("plugins").String(), HomePath)
	cfg.ProvisioningPath = makeAbsolute(iniFile.Section("paths").Key("provisioning").String(), HomePath)
	server := iniFile.Section("server")
	AppUrl, AppSubUrl = parseAppUrlAndSubUrl(server)

	Protocol = HTTP
	if server.Key("protocol").MustString("http") == "https" {
		Protocol = HTTPS
		CertFile = server.Key("cert_file").String()
		KeyFile = server.Key("cert_key").String()
	}
	if server.Key("protocol").MustString("http") == "socket" {
		Protocol = SOCKET
		SocketPath = server.Key("socket").String()
	}

	Domain = server.Key("domain").MustString("localhost")
	HttpAddr = server.Key("http_addr").MustString(DEFAULT_HTTP_ADDR)
	HttpPort = server.Key("http_port").MustString("3000")
	RouterLogging = server.Key("router_logging").MustBool(false)



	// read data proxy settings
	dataproxy := iniFile.Section("dataproxy")
	DataProxyLogging = dataproxy.Key("logging").MustBool(false)



	// read dashboard settings
	dashboards := iniFile.Section("dashboards")
	DashboardVersionsToKeep = dashboards.Key("versions_to_keep").MustInt(20)
	users := iniFile.Section("users")

	//cfg.readSessionConfig()
	cfg.readSmtpSettings()
	//cfg.readQuotaSettings()

	VerifyEmailEnabled = users.Key("verify_email_enabled").MustBool(false)
	if VerifyEmailEnabled && !cfg.Smtp.Enabled {
		//log.Warn("require_email_validation is enabled but smtp is disabled")
	}

	return nil
}


func makeAbsolute(path string, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}


func parseAppUrlAndSubUrl(section *ini.Section) (string, string) {
	appUrl := section.Key("root_url").MustString("http://localhost:3000/")
	if appUrl[len(appUrl)-1] != '/' {
		appUrl += "/"
	}

	// Check if has app suburl.
	url, err := url.Parse(appUrl)
	if err != nil {
		//log.Fatal(4, "Invalid root_url(%s): %s", appUrl, err)
	}
	appSubUrl := strings.TrimSuffix(url.Path, "/")

	return appUrl, appSubUrl
}

func loadConfiguration(args *CommandLineArgs) (*ini.File, error) {
	var err error

	// load config defaults
	defaultConfigFile := path.Join(HomePath, "conf/defaults.ini")
	configFiles = append(configFiles, defaultConfigFile)

	// check if config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		fmt.Println("Merkur-server Init Failed: Could not find config defaults, make sure homepath command line parameter is set or working directory is homepath")
		os.Exit(1)
	}

	// load defaults
	parsedFile, err := ini.Load(defaultConfigFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to parse defaults.ini, %v", err))
		os.Exit(1)
		return nil, err
	}

	parsedFile.BlockMode = false

	// command line props
	commandLineProps := getCommandLineProperties(args.Args)
	// load default overrides
	applyCommandLineDefaultProperties(commandLineProps, parsedFile)

	// load specified config file
	err = loadSpecifedConfigFile(args.Config, parsedFile)
	if err != nil {
		initLogging(parsedFile)

	}

	// apply environment overrides
	err = applyEnvVariableOverrides(parsedFile)
	if err != nil {
		return nil, err
	}

	// apply command line overrides
	applyCommandLineProperties(commandLineProps, parsedFile)

	// evaluate config values containing environment variables
	evalConfigValues(parsedFile)

	// update data path and logging config
	DataPath = makeAbsolute(parsedFile.Section("paths").Key("data").String(), HomePath)
	initLogging(parsedFile)

	return parsedFile, err
}

func initLogging(file *ini.File) {
	// split on comma
	LogModes = strings.Split(file.Section("log").Key("mode").MustString("console"), ",")
	// also try space
	if len(LogModes) == 1 {
		LogModes = strings.Split(file.Section("log").Key("mode").MustString("console"), " ")
	}
	LogsPath = makeAbsolute(file.Section("paths").Key("logs").String(), HomePath)
	log.ReadLoggingConfig(LogModes, LogsPath, file)
}


func evalConfigValues(file *ini.File) {
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			key.SetValue(evalEnvVarExpression(key.Value()))
		}
	}
}


func evalEnvVarExpression(value string) string {
	regex := regexp.MustCompile(`\${(\w+)}`)
	return regex.ReplaceAllStringFunc(value, func(envVar string) string {
		envVar = strings.TrimPrefix(envVar, "${")
		envVar = strings.TrimSuffix(envVar, "}")
		envValue := os.Getenv(envVar)

		// if env variable is hostname and it is empty use os.Hostname as default
		if envVar == "HOSTNAME" && envValue == "" {
			envValue, _ = os.Hostname()
		}

		return envValue
	})
}
func applyEnvVariableOverrides(file *ini.File) error {
	appliedEnvOverrides = make([]string, 0)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			sectionName := strings.ToUpper(strings.Replace(section.Name(), ".", "_", -1))
			keyName := strings.ToUpper(strings.Replace(key.Name(), ".", "_", -1))
			envKey := fmt.Sprintf("GF_%s_%s", sectionName, keyName)
			envValue := os.Getenv(envKey)

			if len(envValue) > 0 {
				key.SetValue(envValue)
				if shouldRedactKey(envKey) {
					envValue = "*********"
				}
				if shouldRedactURLKey(envKey) {
					u, err := url.Parse(envValue)
					if err != nil {
						return fmt.Errorf("could not parse environment variable. key: %s, value: %s. error: %v", envKey, envValue, err)
					}
					ui := u.User
					if ui != nil {
						_, exists := ui.Password()
						if exists {
							u.User = url.UserPassword(ui.Username(), "-redacted-")
							envValue = u.String()
						}
					}
				}
				appliedEnvOverrides = append(appliedEnvOverrides, fmt.Sprintf("%s=%s", envKey, envValue))
			}
		}
	}

	return nil
}
func shouldRedactURLKey(s string) bool {
	uppercased := strings.ToUpper(s)
	return strings.Contains(uppercased, "DATABASE_URL")
}
func applyCommandLineDefaultProperties(props map[string]string, file *ini.File) {
	appliedCommandLineProperties = make([]string, 0)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			keyString := fmt.Sprintf("default.%s.%s", section.Name(), key.Name())
			value, exists := props[keyString]
			if exists {
				key.SetValue(value)
				if shouldRedactKey(keyString) {
					value = "*********"
				}
				appliedCommandLineProperties = append(appliedCommandLineProperties, fmt.Sprintf("%s=%s", keyString, value))
			}
		}
	}
}

func applyCommandLineProperties(props map[string]string, file *ini.File) {
	for _, section := range file.Sections() {
		sectionName := section.Name() + "."
		if section.Name() == ini.DEFAULT_SECTION {
			sectionName = ""
		}
		for _, key := range section.Keys() {
			keyString := sectionName + key.Name()
			value, exists := props[keyString]
			if exists {
				appliedCommandLineProperties = append(appliedCommandLineProperties, fmt.Sprintf("%s=%s", keyString, value))
				key.SetValue(value)
			}
		}
	}
}


func getCommandLineProperties(args []string) map[string]string {
	props := make(map[string]string)

	for _, arg := range args {
		if !strings.HasPrefix(arg, "cfg:") {
			continue
		}

		trimmed := strings.TrimPrefix(arg, "cfg:")
		parts := strings.Split(trimmed, "=")
		if len(parts) != 2 {
			///log.Fatal(3, "Invalid command line argument", arg)
			return nil
		}

		props[parts[0]] = parts[1]
	}
	return props
}
func shouldRedactKey(s string) bool {
	uppercased := strings.ToUpper(s)
	return strings.Contains(uppercased, "PASSWORD") || strings.Contains(uppercased, "SECRET") || strings.Contains(uppercased, "PROVIDER_CONFIG")
}

func loadSpecifedConfigFile(configFile string, masterFile *ini.File) error {
	if configFile == "" {
		configFile = filepath.Join(HomePath, CustomInitPath)
		// return without error if custom file does not exist
		if !pathExists(configFile) {
			return nil
		}
	}

	userConfig, err := ini.Load(configFile)
	if err != nil {
		return fmt.Errorf("Failed to parse %v, %v", configFile, err)
	}

	userConfig.BlockMode = false

	for _, section := range userConfig.Sections() {
		for _, key := range section.Keys() {
			if key.Value() == "" {
				continue
			}

			defaultSec, err := masterFile.GetSection(section.Name())
			if err != nil {
				defaultSec, _ = masterFile.NewSection(section.Name())
			}
			defaultKey, err := defaultSec.GetKey(key.Name())
			if err != nil {
				defaultKey, _ = defaultSec.NewKey(key.Name(), key.Value())
			}
			defaultKey.SetValue(key.Value())
		}
	}

	configFiles = append(configFiles, configFile)
	return nil
}

