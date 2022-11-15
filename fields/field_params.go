package fields

type ParamInt struct {
	IsSet bool
	Value int
}
func NewParamInt(value int) ParamInt{
	return ParamInt{IsSet: true, Value: value}
}

type ParamInt64 struct {
	IsSet bool
	Value int64
}
func NewParamInt64(value int64) ParamInt64{
	return ParamInt64{IsSet: true, Value: value}
}

type ParamFloat struct {
	IsSet bool
	Value float64
}
func NewParamFloat(value float64) ParamFloat{
	return ParamFloat{IsSet: true, Value: value}
}

type ParamBool struct {
	IsSet bool
	Value bool
}
func NewParamBool(value bool) ParamBool{
	return ParamBool{IsSet: true, Value: value}
}

