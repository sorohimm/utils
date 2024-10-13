package cfg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type yamlCfg struct {
	Test struct {
		Http struct {
			URL     string `yaml:"url" env:"HTTP_URL_NOT_VALID" validate:"required,url"`
			Port    int    `yaml:"port" env:"HTTP_PORT_NOT_VALID" validate:"required,number"`
			Timeout struct {
				Idle time.Duration `yaml:"idle" validate:"required"`
			} `yaml:"timeout"`
		} `yaml:"http"`
	} `yaml:"test"`
}

type yamlEnvCfg struct {
	Test struct {
		Http struct {
			Url string `yaml:"url" env:"HTTP_URL" validate:"required,url"`
		} `yaml:"http"`
	} `yaml:"test"`
}

type notValidCfg struct {
	Test struct {
		Url string `yaml:"urlNotValid" env:"AUTH_URL_NOT_VALID" validate:"required,url"`
	} `yaml:"test"`
}

func TestCfg(t *testing.T) {
	err := Setup(&SetupParams{
		Prefix:                 "",
		DevPath:                "test.yaml",
		LookupDepth:            LookupDepthDefault,
		TargetEnvFileExtension: ".env.test",
	})
	if err != nil {
		return
	}

	t.Run("testYaml", func(t *testing.T) {
		testYaml(t)
	})
	t.Run("testYamlEnv", func(t *testing.T) {
		testYamlEnv(t)
	})
	t.Run("testDotEnvDoesNotOverrideEnv", func(t *testing.T) {
		testDotEnvDoesNotOverrideEnv(t)
	})
	t.Run("testValidation", func(t *testing.T) {
		testValidation(t)
	})
}

func testYaml(t *testing.T) {
	t.Helper()
	c := yamlCfg{}
	err := Load(&c)
	require.NoError(t, err)
	require.Equal(t, "https://exapmle.com", c.Test.Http.URL)
}

func testYamlEnv(t *testing.T) {
	t.Helper()
	c := yamlEnvCfg{}
	err := Load(&c)
	require.NoError(t, err)
	require.Equal(t, "https://exapmle.com", c.Test.Http.Url)
	t.Setenv(Var("HTTP_URL"), "http://a.b.c")
	err = Load(&c)
	require.NoError(t, err)
	require.Equal(t, "http://a.b.c", c.Test.Http.Url)
}

func testDotEnvDoesNotOverrideEnv(t *testing.T) {
	t.Helper()
	t.Setenv(Var("HTTP_URL"), "http://1.2.c")
	c := yamlEnvCfg{}
	err := Load(&c)
	require.NoError(t, err)
	require.Equal(t, "http://1.2.c", c.Test.Http.Url)
}

func testValidation(t *testing.T) {
	t.Helper()
	c := notValidCfg{}
	err := Load(&c)
	require.Error(t, err)
	require.ErrorContains(t, err, "Field validation for 'Url' failed on the 'required' tag")
}
