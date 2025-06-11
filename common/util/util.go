package util

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&dest); err != nil {
		logrus.Errorf("failed to unmarshal config: %v", err)
		return err
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	if err := v.Unmarshal(&env); err != nil {
		logrus.Errorf("failed to unmarshal config: %v", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			return fmt.Errorf("unsupported env type for key %s: %s", k, valOf.Kind())
		}

		err := os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil
}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")

	if err := v.AddRemoteProvider("consul", endPoint, path); err != nil {
		logrus.Errorf("failed to add consul remote provider: %v", err)
		return err
	}

	if err := v.ReadRemoteConfig(); err != nil {
		logrus.Errorf("failed to read remote config: %v", err)
		return err
	}

	if err := v.Unmarshal(&dest); err != nil {
		logrus.Errorf("failed to unmarshal remote config: %v", err)
		return err
	}

	if err := SetEnvFromConsulKV(v); err != nil {
		logrus.Errorf("failed to set env from consul: %v", err)
		return err
	}

	return nil
}
