package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

const (
	priorityFile = iota
	priorityEnvVars
	priorityOsArgs
	topPriority
)

type Config struct {
	srvAddr       string
	dbDSN         string
	sessionLifeTime time.Duration
}

func New(opts ...configOption) *Config {
	configOptions := setConfigOptions(opts)

	configs := make([]*Config, topPriority)

	configs[priorityEnvVars] = fromEnv(configOptions.envVars)

	if !configOptions.ignoreOsArgs {
		configs[priorityOsArgs] = fromArgs(configOptions.osArgs, &configOptions.filePath)
	}

	if !configOptions.ignoreCfgFile && configOptions.filePath != "" {
		configs[priorityFile] = fromFile(configOptions.filePath)
	}

	return fillCfg(configs)
}

// fillCfg fills Config from passed collection of pConfig
// from different sources considering priority.
func fillCfg(configs []*Config) *Config {
	cfg := Config{}

	for _, pCfg := range configs {
		if pCfg == nil {
			continue
		}

		if pCfg.srvAddr != "" {
			cfg.srvAddr = pCfg.srvAddr
		}
		if pCfg.dbDSN != "" {
			cfg.dbDSN = pCfg.dbDSN
		}
	}

	return cfg.setDefaults()
}

func (c Config) SrvAddr() string {
	return c.srvAddr
}

func (c Config) DBDSN() string {
	return c.dbDSN
}

func (c Config) SessionLifetime() time.Duration {
	return c.sessionLifeTime
}

func (c *Config) setDefaults() *Config {
	if c.srvAddr == "" {
		c.srvAddr = "localhost:8080"
	}

	if c.sessionLifeTime == time.Duration(0) {
		c.sessionLifeTime = time.Minute * 30
	}

	return c
}


func fromEnv(envVars map[string]string) *Config {
	pc := Config{}

	if v := envVars["SERVER_ADDRESS"]; v != "" {
		pc.srvAddr = v
	}
	if v := envVars["DATABASE_DSN"]; v != "" {
		pc.dbDSN = v
	}
	if v := envVars["SESSION_LIFETIME"]; v != "" {
		pc.dbDSN = v
	}

	return &pc
}

func fromArgs(osArgs []string, filePath *string) *Config {
	pc := Config{}
	fs := flag.NewFlagSet("myFS", flag.ContinueOnError)
	if !fs.Parsed() {
		fs.StringVar(&pc.srvAddr, "a", "", "the shortener service address")
		fs.StringVar(&pc.dbDSN, "d", "", "db connection path")
		fs.DurationVar(&pc.sessionLifeTime, "s", time.Duration(0), "session lifetime")

		fs.Parse(osArgs)
	}

	return &pc
}

func fromFile(filePath string) *Config {
	pc := Config{}
	fileData, err := parseFile(filePath)

	if err != nil {
		log.Printf("failed to parse config file %v: %v", filePath, err)
		return &pc
	}

	pc.dbDSN = fileData.DatabaseDsn
	pc.srvAddr = fileData.ServerAddress
	pc.sessionLifeTime = time.Minute * time.Duration(fileData.SessionLifeTime)

	return &pc
}

type fileStruct struct {
	ServerAddress   string `json:"server_address"`
	DatabaseDsn     string `json:"database_dsn"`
	SessionLifeTime int `json:"session_lifetime"` // in minutes
}

func parseFile(p string) (*fileStruct, error) {
	f, err := os.OpenFile(p, os.O_RDONLY, 0o777)
	if err != nil {
		return nil, err
	}

	data := fileStruct{}

	dec := json.NewDecoder(f)
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

