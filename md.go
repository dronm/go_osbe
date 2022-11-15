package osbe

import(
	"time"
	
	"osbe/model"
//	"osbe/fields"
)

const (
	DEF_FIELD_SEP = "@@"
)

type VersionType struct {
	DateOpen time.Time
	DateClose time.Time
	Value string
}

type ModelMDCollection map[string]*model.ModelMD

type Metadata struct {
	Debug bool
	Owner string
	DataSchema string
	Version VersionType
	Controllers ControllerCollection
	Models ModelMDCollection
	Enums EnumCollection
	Constants ConstantCollection
}

func NewMetadata() *Metadata{
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

