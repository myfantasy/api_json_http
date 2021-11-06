package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/myfantasy/authentication/sat"
	"github.com/myfantasy/authorization/saz"
	"github.com/myfantasy/compress"
	"github.com/myfantasy/mft"
	"github.com/myfantasy/storage"
	log "github.com/sirupsen/logrus"
)

type settingsParcer struct {
	SettingsMountPoint string                `json:"settings_mount_point"`
	LogLevel           string                `json:"log_level"`
	ZipLevel           int                   `json:"zip_level"`
	Storage            storage.GeneratorInfo `json:"storage"`

	SimpleAuthenticationFilename string `json:"simple_authentication_filename"`
	SimpleAuthorizationFilename  string `json:"simple_authorization_filename"`

	DefaultWriteTimeout time.Duration `json:"default_write_timeout"`

	SettingsStorage storage.Storage `json:"-"`
}

var settings settingsParcer

func loadSettings() {
	settingsData, er0 := ioutil.ReadFile(*fSettingsFile)
	if er0 != nil {
		log.Fatalf("File read fail: %v", er0)
	}

	settings.LogLevel = "debug"
	settings.ZipLevel = 7

	settings.SimpleAuthenticationFilename = "authentication.json"
	settings.SimpleAuthorizationFilename = "authorization.json"

	settings.DefaultWriteTimeout = 5 * time.Second
	er0 = json.Unmarshal(settingsData, &settings)
	if er0 != nil {
		log.Fatalf("File parce fail: %v", er0)
	}

	compressor = compress.GeneratorCreate(settings.ZipLevel)

	fileStorage = storage.CreateGenerator(settings.Storage, compressor)

	var err *mft.Error

	ctx, cancel := context.WithTimeout(context.Background(), settings.DefaultWriteTimeout)

	settings.SettingsStorage, err = fileStorage.Create(
		ctx, settings.SettingsMountPoint, "")
	cancel()

	if err != nil {
		log.Fatalf("Create settings file storages error: %v", err)
	}

	autht = &sat.SimpleAuthenticationChecker{}
	autht.SaveToContextDuration = settings.DefaultWriteTimeout
	autht.SaveToFileNameValue = settings.SimpleAuthenticationFilename
	autht.SaveToStorageValue = settings.SettingsStorage

	err = autht.Load()
	if err != nil {
		log.Fatalf("SimpleAuthenticationChecker load error: %v", err)
	}

	authz = &saz.SimplePermissionChecker{}
	authz.SaveToContextDuration = settings.DefaultWriteTimeout
	authz.SaveToFileNameValue = settings.SimpleAuthorizationFilename
	authz.SaveToStorageValue = settings.SettingsStorage

	err = authz.Load()
	if err != nil {
		log.Fatalf("SimplePermissionChecker load error: %v", err)
	}
}
