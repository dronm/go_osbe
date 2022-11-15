package fields

type FieldDataType byte

const(
	FIELD_TYPE_BOOL FieldDataType = iota
	FIELD_TYPE_CHAR
	FIELD_TYPE_STRING
	FIELD_TYPE_INT
	FIELD_TYPE_DATE
	FIELD_TYPE_TIME
	FIELD_TYPE_TIMETZ
	FIELD_TYPE_DATETIME
	FIELD_TYPE_DATETIMETZ
	FIELD_TYPE_FLOAT
	FIELD_TYPE_TEXT
	FIELD_TYPE_ENUM
	FIELD_TYPE_PASSWORD
	FIELD_TYPE_INTERVAL
	FIELD_TYPE_JSON
	FIELD_TYPE_JSONB
	FIELD_TYPE_ARRAY
	FIELD_TYPE_BYTEA
	FIELD_TYPE_XML
	FIELD_TYPE_GEOMPOLYGON
	FIELD_TYPE_GEOMPOINT	
)
	
//Base metadata field
type Field struct {
	Id string
	PrimaryKey bool
	AutoInc bool
	Required bool
	SysCol bool
	DbRequired bool
	Display bool
	Alias string
	Descr string
	Length int
	DefOrder ParamBool
	DefOrderIndex byte //to preserve order
	//DefaultValue interface{}
	RefTable string
	RefField string
	Unique bool
	EnumId string
	Precision int
	RegFieldType string
	RetAfterInsert bool
	NoValueOnCopy bool
	OrderInList byte //to preserve order for getting comma separated list for select queries
	Encrypted bool
}
func (f *Field) GetId() string {
	return f.Id
}
func (f *Field) SetId(v string) {
	f.Id = v
}

func (f *Field) GetRequired() bool {
	return f.Required
}
func (f *Field) SetRequired(v bool) {
	f.Required = v
}

func (f *Field) GetAlias() string {
	return f.Alias
}
func (f *Field) SetAlias(v string) {
	f.Alias = v
}
func (f *Field) GetDefOrder() ParamBool {
	return f.DefOrder
}
func (f *Field) SetDefOrder(v ParamBool) {
	f.DefOrder = v
}

func (f *Field) GetDescr() string {
	if f.Descr != "" {
		return f.Descr
	}else{
		if f.Alias != "" {
			return f.Alias
		}else{
			return f.Id
		}		
	}	
}

func (f *Field) SetDescr(v string) {
	f.Descr = v
}

/*func (f *Field) GetDataType() FieldDataType {
	return f.DataType
}*/
func (f *Field) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *Field) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}

func (f *Field) GetAutoInc() bool {
	return f.AutoInc
}
func (f *Field) GetSysCol() bool {
	return f.SysCol
}
func (f *Field) GetDbRequired() bool {
	return f.DbRequired
}
func (f *Field) GetDisplay() bool {
	return f.Display
}
func (f *Field) GetLength() int {
	return f.Length
}
func (f *Field) GetRefTable() string {
	return f.RefTable
}
func (f *Field) GetRefField() string {
	return f.RefField
}
func (f *Field) GetUnique() bool {
	return f.Unique
}
func (f *Field) GetEnumId() string {
	return f.EnumId
}
func (f *Field) GetPrecision() int {
	return f.Precision
}
func (f *Field) GetRegFieldType() string {
	return f.RegFieldType
}
func (f *Field) GetRetAfterInsert() bool {
	return f.RetAfterInsert
}
func (f *Field) GetNoValueOnCopy() bool {
	return f.NoValueOnCopy
}
func (f *Field) GetOrderInList() byte {
	return f.OrderInList
}
func (f *Field) SetOrderInList(v byte) {
	f.OrderInList = v
}

func (f *Field) GetEncrypted() bool {
	return f.Encrypted
}
func (f *Field) SetEncrypted(v bool) {
	f.Encrypted = v
}
func (f *Field) GetDefOrderIndex() byte {
	return f.DefOrderIndex
}
func (f *Field) SetDefOrderIndex(v byte) {
	f.DefOrderIndex = v
}


//Base interface
type Fielder interface {
	GetId() string
	SetId(string)		
	GetAlias() string
	SetAlias(string)
	GetRequired() bool
	SetRequired(bool)
	GetDescr() string
	SetDescr(string)
	GetDataType() FieldDataType
	GetPrimaryKey() bool
	SetPrimaryKey(bool)
	GetAutoInc() bool
	GetSysCol() bool
	GetDbRequired() bool
	GetDisplay() bool
	GetRegFieldType() string
	GetRetAfterInsert() bool
	GetNoValueOnCopy() bool	
	GetDefOrder() ParamBool
	SetDefOrder(v ParamBool)
	GetOrderInList() byte
	SetOrderInList(byte)
	GetDefOrderIndex() byte
	SetDefOrderIndex(byte)
	SetEncrypted(bool)
	GetEncrypted() bool
}

type FieldCollection map[string]Fielder

