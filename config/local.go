package config

import (
	"logsnap/constants"
	"logsnap/version"
	"time"

	"github.com/sirupsen/logrus"
)

type LocalConfig struct {
	RemoteConfigLastCheck string
}

func NewLocalConfig() *LocalConfig {
	return &LocalConfig{}
}

func (c *LocalConfig) GetVersion() string {
	return version.GetVersion()
}

func (c *LocalConfig) GetRemoteConfigEnabled() bool {
	return false
}

func (c *LocalConfig) GetUploadConfigURL() string {
	return constants.UploadConfigURL
}

func (c *LocalConfig) GetDownloadConfigURL() string {
	logrus.Debugf("DownloadConfigURL: %s", constants.DownloadConfigURL)
	return constants.DownloadConfigURL
}

func (c *LocalConfig) GetRemoteConfigInterval() int {
	return 36000
}

func (c *LocalConfig) GetRemoteConfigLastCheck() string {
	return time.Now().Format(time.RFC3339)
}

func (c *LocalConfig) SetRemoteConfigLastCheck(lastCheck string) {
	c.RemoteConfigLastCheck = lastCheck
}
