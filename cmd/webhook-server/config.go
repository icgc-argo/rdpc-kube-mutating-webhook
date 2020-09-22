/*
 * Copyright (c) 2020. Ontario Institute for Cancer Research
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mcuadros/go-defaults.v1"
	"gopkg.in/yaml.v2"
	"os"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(configFilePath string, cfg *Config) {
	f, err := os.Open(configFilePath)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func parseConfig() Config {
	var cfg = Config{}
	defaults.SetDefaults(&cfg)
	if len(os.Args) > 1 {
		var configFile = os.Args[1]
		readFile(configFile, &cfg)
	}
	readEnv(&cfg)
	var validate = validator.New()
	var err = validate.Struct(&cfg)

	if err != nil {
		panic(err.Error())
	}
	return cfg
}

type Config struct {
	Server struct {
		Port string `default:"8080", validate:"required", yaml:"port", envconfig:"SERVER_PORT"`
		SSL  struct {
			Enable   bool   `default:"false", validate:"required", yaml:"enable", envconfig:"SERVER_SSL_ENABLE"`
			CertPath string `yaml:"certPath", envconfig:"SERVER_SSL_CERTPATH"`
			KeyPath  string `yaml:"keyPath", envconfig:"SERVER_SSL_KEYPATH"`
		} `yaml:"ssl"`
	} `yaml:"server"`
	App struct {
		Debug bool `default:"false", validate:"required",yaml:"debug",envconfig:"APP_DEBUG"`
		OverrideVolumeCollisions bool   `default:"false", validate:"required", yaml:"overrideVolumeCollisions", envconfig:"APP_OVERRIDEVOLUMECOLLISIONS"`
		TargetContainerName      string `validate:"required", yaml:"targetContainerName", envconfig:"APP_TARGETCONTAINERNAME"`
		EmptyDir                 struct {
			VolumeName string `validate:"required", yaml:"volumeName", envconfig:"APP_EMPTYDIR_VOLUMENAME"`
			MountPath  string `validate:"required", yaml:"mountPath", envconfig:"APP_EMPTYDIR_MOUNTPATH"`
		} `yaml:"emptydir"`
	} `yaml:"app"`
}
