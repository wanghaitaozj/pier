package repo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// DefaultPathName is the default config dir name
	DefaultPathName = ".pier"

	// DefaultPathRoot is the path to the default config dir location.
	DefaultPathRoot = "~/" + DefaultPathName

	// EnvDir is the environment variable used to change the path root.
	EnvDir = "PIER_PATH"

	// ConfigName is config name
	ConfigName = "pier.toml"

	// config path
	ConfigPath = "../../config"

	// KeyName
	KeyName = "key.json"

	// NodeKeyName
	NodeKeyName = "node.priv"

	// KeyPassword
	KeyPassword = "bitxhub"

	// API name
	APIName = "api"
)

var RootPath string

// Initialize creates .pier path and necessary configuration,
// account file and so on.
func Initialize(repoRoot string) error {
	if _, err := os.Stat(repoRoot); os.IsNotExist(err) {
		err := os.MkdirAll(repoRoot, 0755)
		if err != nil {
			return err
		}
	}

	box := packr.NewBox(ConfigPath)

	privKey, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return fmt.Errorf("create private key error: %s", err)
	}

	keyPath := filepath.Join(repoRoot, KeyName)
	if err := asym.StorePrivateKey(privKey, keyPath, KeyPassword); err != nil {
		return fmt.Errorf("persist key: %s", err)
	}

	if err := box.Walk(func(s string, file packd.File) error {
		p := filepath.Join(repoRoot, s)
		dir := filepath.Dir(p)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(p, []byte(file.String()), 0644)
	}); err != nil {
		return err
	}

	return nil
}

// InitConfig initialize configuration
func InitConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvDir)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})

	return nil
}

// PathRoot returns root path (default .pier)
func PathRoot() (string, error) {
	if RootPath != "" {
		return RootPath, nil
	}
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(DefaultPathRoot)
	}

	return dir, err
}

// SetPath sets global config path
func SetPath(root string) {
	RootPath = root
}

// PathRootWithDefault gets current config path with default value
func PathRootWithDefault(path string) (string, error) {
	if len(path) == 0 {
		return PathRoot()
	}

	return path, nil
}

// PluginPath returns the plugin path
func PluginPath() (string, error) {
	repoRoot, err := PathRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(repoRoot, "plugins"), nil
}

func KeyPath(repoRoot string) string {
	return filepath.Join(repoRoot, KeyName)
}

func NodeKeyPath(repoRoot string) string {
	return filepath.Join(repoRoot, NodeKeyName)
}

// LoadPrivateKey loads private key from config path
func LoadPrivateKey(repoRoot string) (crypto.PrivateKey, error) {
	keyPath := filepath.Join(repoRoot, KeyName)

	return asym.RestorePrivateKey(keyPath, "bitxhub")
}

func GetAPI(repoRoot string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(repoRoot, APIName))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
