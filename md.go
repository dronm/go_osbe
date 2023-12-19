package osbe

import(
	"time"
	
	"osbe/model"
)

// This file describes MD object. MD object holds descriptions of
// application objects: controllers, models, constants, enums.
// Actually it holds collections of these objects.

const (
	DEF_FIELD_SEP = "@@" //default field separator used by clients
)

// VersionType structures holds metadata version information.
type VersionType struct {
	DateOpen time.Time	// open for modification timestamp
	DateClose time.Time	// close version timestamp
	Value string		// version value
}

// ModelMDCollection
type ModelMDCollection map[string]*model.ModelMD

// Metadata is the metadata structure.
// Maps metadata.xml file.
type Metadata struct {
	Debug bool
	Owner string				// holds owner information
	DataSchema string
	Version VersionType
	Controllers ControllerCollection
	Models ModelMDCollection
	Enums EnumCollection
	Constants ConstantCollection
}
// NewMetadata creates new Metadata object
func NewMetadata() *Metadata {
	return &Metadata{Controllers: make(ControllerCollection),
			Constants: make(ConstantCollection),
			Models: make(ModelMDCollection),
	}
}

type LocaleID string
const (
	LOCALE_RU LocaleID = "ru"
	LOCALE_EN LocaleID = "en"	
)

