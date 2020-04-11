package flag

import "strconv"

type FlagType int

func (F FlagType) String() string {
	switch F {
	case StringFlag:
		return "String"
	case IntegerFlag:
		return "Integer"
	case FloatFlag:
		return "Float"
	case JsonFlag:
		return "JSON"
	default:
		return ""
	}
}

const (
	StringFlag FlagType = iota
	BooleanFlag
	IntegerFlag
	FloatFlag
	JsonFlag
)

type FeatureFlag struct {
	Namespace string
	FlagName  string
	Value     interface{} // Raw value contained in flag
	Type      FlagType
	Refresh   int // Rate for client to refresh flag
}

func NewFlag() *FeatureFlag {
	return &FeatureFlag{}
}

func (fl *FeatureFlag) String() string {
	switch fl.Type {
	case StringFlag:
		return fl.Value.(string)
	case IntegerFlag:
		return strconv.Itoa(fl.Value.(int))
	case BooleanFlag:
		return strconv.FormatBool(fl.Value.(bool))
	default:
		return ""
	}
}
