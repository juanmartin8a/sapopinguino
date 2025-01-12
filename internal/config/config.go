package config

// READS THE FILES FROM ./config

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	awsutils "sapopinguino/internal/aws"

	"github.com/spf13/viper"
)

type config struct {
	AppEnv   string
	Database struct {
		DSN string
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
	config := &C

	e := appEnv(option)

	if e == Production {
		setProd()
	} else {
		setDev()
	}

	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

    err := setCValues(&C.Database.DSN, &C.OpenAI.Key)
	if err != nil {
		log.Fatal(err)
	}
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

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func setDev() {
	viper.AddConfigPath(
		// filepath.Join(rootDir(), "config"),
        "/config",
	)
	viper.SetConfigName("config.dev")
}

func setProd() {
	viper.AddConfigPath(
		// filepath.Join(rootDir(), "config"),
        "/config",
	)
	viper.SetConfigName("config.prod")
}

func setCValues(params ...*string) error {
	for _, paramName := range params {
		paramBytes, err := awsutils.GetSecretString(
			*paramName,
		)
		if err != nil {
			log.Printf(
				"Error while getting param %s: %s",
				*paramName,
				err,
			)
			return err
		}

		param := string(paramBytes)

		*paramName = param
	}

	return nil
}
