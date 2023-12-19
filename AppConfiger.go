package osbe

import (
	"osbe/config"
)

// AppConfiger is an application configuration
// interface. Realization is held by Application object.
type AppConfiger interface {
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
	GetDebugQueries() bool
}


