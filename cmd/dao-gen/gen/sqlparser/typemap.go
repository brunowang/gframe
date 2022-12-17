package sqlparser

import "github.com/pingcap/tidb/parser/mysql"

var typeMap = map[byte]string{
	mysql.TypeTiny:       "int8",
	mysql.TypeShort:      "int16",
	mysql.TypeLong:       "int32",
	mysql.TypeFloat:      "float32",
	mysql.TypeDouble:     "float64",
	mysql.TypeTimestamp:  "time.Time",
	mysql.TypeLonglong:   "int64",
	mysql.TypeInt24:      "int32",
	mysql.TypeDate:       "time.Time",
	mysql.TypeDuration:   "time.Time",
	mysql.TypeDatetime:   "time.Time",
	mysql.TypeYear:       "time.Time",
	mysql.TypeNewDate:    "time.Time",
	mysql.TypeVarchar:    "string",
	mysql.TypeBit:        "int8",
	mysql.TypeJSON:       "map[string]interface{}",
	mysql.TypeTinyBlob:   "string",
	mysql.TypeMediumBlob: "string",
	mysql.TypeLongBlob:   "string",
	mysql.TypeBlob:       "string",
	mysql.TypeVarString:  "string",
	mysql.TypeString:     "string",
}
