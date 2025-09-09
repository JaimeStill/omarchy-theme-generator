package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	ConfigFile   = "omarchy-theme-gen"
	ConfigDir    = "omarchy"
	ConfigEnv    = "XDG_CONFIG_HOME"
	ConfigFormat = "json"
	EnvPrefix    = "OMARCHY_THEME_GEN"
	SystemDir    = "/etc"
)

func Load() (*Settings, error) {
	v := viper.New()

	setDefaults(v)

	// Check for explicit config file path first
	if configFile := os.Getenv("OMARCHY_CONFIG"); configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			// If explicit config file is specified but can't be read, return error
			return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
		}
	} else {
		// Use default config search paths
		v.SetConfigName(ConfigFile)
		v.SetConfigType(ConfigFormat)

		v.AddConfigPath(filepath.Join(SystemDir, ConfigDir))

		if xdgConfig := os.Getenv(ConfigEnv); xdgConfig != "" {
			v.AddConfigPath(filepath.Join(xdgConfig, ConfigDir))
		}

		v.AddConfigPath(".")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config: %w", err)
			}
		}
	}

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var settings Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &settings, nil
}

func LoadWithViper(v *viper.Viper) (*Settings, error) {
	var settings Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	return &settings, nil
}

func SaveToFile(settings *Settings, path string) error {
	v := viper.New()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return v.WriteConfigAs(path)
}

func GetUserConfigPath() string {
	xdgConfig := os.Getenv(ConfigEnv)
	return filepath.Join(xdgConfig, ConfigDir, ConfigFile+"."+ConfigFormat)
}

func GetSystemConfigPath() string {
	return filepath.Join(SystemDir, ConfigDir, ConfigFile+"."+ConfigFormat)
}
