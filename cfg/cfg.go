package cfg

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const (
	devEnv             = "dev"
	stageEnv           = "stage"
	prodEnv            = "prod"
	LookupDepthDefault = 6
)

var (
	validate   = validator.New()
	cfg        *config
	onceConfig sync.Once
)

type config struct {
	prefix string
	file   []byte
}

type SetupParams struct {
	Prefix                 string
	DevPath                string
	StagePath              string
	ProdPath               string
	LookupDepth            int
	TargetEnvFileExtension string
}

type configLoader struct {
	params *SetupParams
	cfg    *config
	err    error
}

func NewPrefix(p string) string {
	return p + "_"
}

// Setup is a function that sets up the configuration for the application.
//
// It takes a pointer to a SetupParams struct as a parameter and returns an error.
// The SetupParams struct contains the following fields:
// - Prefix: a string specifying the prefix for environment variables
// - DevPath: a string specifying the path to the development config file
// - StagePath: a string specifying the path to the staging config file
// - ProdPath: a string specifying the path to the production config file
// - LookupDepth: an integer specifying the depth for file lookup
// - TargetEnvFileExtension: a string specifying the extension of the target environment file
//
// It returns any error that occurred during the setup process.
//
// Example usage:
//
//	err := Setup(&SetupParams{
//	    Prefix: "",
//	    DevPath: "test.yaml",
//	    LookupDepth: LookupDepthDefault,
//	    TargetEnvFileExtension: ".env.test",
//	})
func Setup(params *SetupParams) error {
	loader := &configLoader{params: params}
	return loader.setupConfig()
}

func (c *configLoader) setupConfig() error {
	c.setDefaultValues()
	c.loadEnvFile()
	c.checkEnvVariable()

	switch os.Getenv(c.params.Prefix + "ENV") {
	case devEnv:
		c.loadConfigFile(c.params.DevPath)
	case stageEnv:
		c.loadConfigFile(c.params.StagePath)
	case prodEnv:
		c.loadConfigFile(c.params.ProdPath)
	}

	return c.err
}

func (c *configLoader) setDefaultValues() {
	onceConfig.Do(func() {
		if c.params.LookupDepth == 0 {
			c.params.LookupDepth = LookupDepthDefault
		}
	})
}

func (c *configLoader) checkEnvVariable() {
	envKey := c.params.Prefix + "ENV"
	envValue := os.Getenv(envKey)
	if c.err == nil && envValue == "" {
		c.err = fmt.Errorf("ENV variable is undefined")
	}
}

func (c *configLoader) loadEnvFile() {
	lookup := NewLookup(c.params.TargetEnvFileExtension, c.params.LookupDepth)
	envFile, _ := lookup.FindFile()
	err := loadEnv(envFile)
	if c.err == nil && err != nil {
		c.err = fmt.Errorf("load %s file error: %w", envFile, err)
	}
}

func (c *configLoader) loadConfigFile(filePath string) {
	if c.err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			c.err = fmt.Errorf("unable to read %s file: %w", filePath, err)
			return
		}
		c.cfg = &config{
			prefix: c.params.Prefix,
			file:   file,
		}
		cfg = c.cfg
	}
}

// Load is a function that loads the configuration data into the destination struct.
//
// Example usage:
// c := Cfg{}
// err := Load(&c)
//
//	if err != nil {
//	    // handle error
//	}
func Load(dst interface{}) error {
	if cfg == nil {
		return fmt.Errorf("config is nil, check setup")
	}
	if err := yaml.Unmarshal(cfg.file, dst); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	opts := env.Options{
		UseFieldNameByDefault: false,
		Prefix:                cfg.prefix,
	}

	if err := env.ParseWithOptions(dst, opts); err != nil {
		return fmt.Errorf("parse env error: %w", err)
	}

	return validateConfig(dst)
}

func validateConfig(dst interface{}) error {
	if err := validate.Struct(dst); err != nil {
		return accumulateError(err)
	}
	return nil
}

func accumulateError(err error) Error {
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		var e Error
		for _, fe := range valErrs {
			e.errors = append(e.errors, fe)
		}
		return e
	}
	return Error{}
}
