package config

// READS THE .yml FILES FROM ./config

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	AppEnv    string
	Websocket struct {
		Endpoint string
	}
	OpenAI struct {
		Key string
	}
}

var C config

const (
	Development = "dev"
	Production  = "prod"
)

type ReadConfigOption struct {
	AppEnv string
}

func ReadConfig(option ReadConfigOption) {
	e := appEnv(option)

	if e == Production {
		setProd()
	} else {
		setDev()
	}

	viper.SetConfigType("yml")

	LoadDotEnv()

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalln(err)
	}

	// err := setCValues(&C.OpenAI.Key)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func appEnv(option ReadConfigOption) string {
	if option.AppEnv != "" {
		return option.AppEnv
	}
	if os.Getenv("APP_ENV") != "" {
		return os.Getenv("APP_ENV")
	}

	return Development
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d) // internal/
}

func setDev() {
	viper.AddConfigPath(
		filepath.Join(RootDir(), "config"),
		// "/config",
	)
	viper.SetConfigName("config.dev")
}

func setProd() {
	viper.AddConfigPath(
		filepath.Join(RootDir(), "config"),
		// "/config",
	)
	viper.SetConfigName("config.prod")
}

// func setCValues(params ...*string) error {
// 	for _, paramName := range params {
// 		paramBytes, err := awsutils.GetSecretString(
// 			*paramName,
// 				)
// 		if err != nil {
// 			log.Printf(
// 				"Error while getting param %s: %s",
// 				*paramName,
// 				err,
// 			)
// 			return err
// 		}
//
// 		param := string(paramBytes)
//
// 		*paramName = param
// 	}
//
// 	return nil
// }
