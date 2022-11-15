package osbe

import (
	"osbe/config"
)

type AppConfiger interface {
	ReadConf(fileName string) error
	WriteConf(fileName string) error
	GetDb() config.DbStorage
	GetWSServer() string
	GetTLSWSServer() string
	GetTLSKey() string
	GetTLSCert() string
	GetAppID() string
	GetLogLevel() string
	GetSession() config.Session
	GetTemplateDir() string
	GetReportErrors() bool
	GetXSLTDir() string
	GetDefaultLocale() string
	GetTechMail() string
	GetAuthor() string
}


