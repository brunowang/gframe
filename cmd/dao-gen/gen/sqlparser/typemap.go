package sqlparser

import "github.com/pingcap/tidb/pkg/parser/mysql"

var typeMap = TypeMap{
	mysql.TypeTiny:       {"int8", "0"},
	mysql.TypeShort:      {"int16", "0"},
	mysql.TypeLong:       {"int32", "0"},
	mysql.TypeFloat:      {"float32", "0.0"},
	mysql.TypeDouble:     {"float64", "0.0"},
	mysql.TypeTimestamp:  {"time.Time", "time.Time{}"},
	mysql.TypeLonglong:   {"int64", "0"},
	mysql.TypeInt24:      {"int32", "0"},
	mysql.TypeDate:       {"time.Time", "time.Time{}"},
	mysql.TypeDuration:   {"time.Time", "time.Time{}"},
	mysql.TypeDatetime:   {"time.Time", "time.Time{}"},
	mysql.TypeYear:       {"time.Time", "time.Time{}"},
	mysql.TypeNewDate:    {"time.Time", "time.Time{}"},
	mysql.TypeVarchar:    {"string", "\"\""},
	mysql.TypeBit:        {"int8", "0"},
	mysql.TypeJSON:       {"map[string]interface{}", "make(map[string]interface{})"},
	mysql.TypeTinyBlob:   {"string", "\"\""},
	mysql.TypeMediumBlob: {"string", "\"\""},
	mysql.TypeLongBlob:   {"string", "\"\""},
	mysql.TypeBlob:       {"string", "\"\""},
	mysql.TypeVarString:  {"string", "\"\""},
	mysql.TypeString:     {"string", "\"\""},
}

type TypeMap map[byte]TypeInfo

func (m TypeMap) GetType(typ byte) string {
	tp, ok := m[typ]
	if ok {
		return tp.Type
	}
	return ""
}

func (m TypeMap) GetZeroVal(typ byte) string {
	tp, ok := m[typ]
	if ok {
		return tp.ZeroVal
	}
	return ""
}

type TypeInfo struct {
	Type    string
	ZeroVal string
}
