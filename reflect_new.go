package goutils

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	typeRegistry = make(map[string]reflect.Type)
)

func RegisterType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.Name()] = t
}

func MakeInstance(name string) interface{} {
	return reflect.New(typeRegistry[name]).Elem().Interface()
}

func MakeInstancePtr(name string) interface{} {
	return reflect.New(typeRegistry[name]).Interface()
}

func StructScan(rows *sql.Rows, t reflect.Type) (o interface{}, err error) {
	cols, _ := rows.Columns()

	var m map[string]interface{}
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err = rows.Scan(columnPointers...); err != nil {
		return
	}

	m = make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		if *val == nil {
			continue
		}

		m[colName] = *val
	}

	v := reflect.New(t).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]

		if item, ok := m[field]; ok {
			//reflect.Typeof(item) can be: Time, int64, float32, float64, []uint8
			if v.Field(i).CanSet() {
				if item != nil {
					switch v.Field(i).Kind() {
					case reflect.String:
						v.Field(i).SetString(string(item.([]byte)))
					case reflect.Float32, reflect.Float64:
						switch fl := item.(type) {
						case float32:
							v.Field(i).SetFloat(float64(fl))
						case float64:
							v.Field(i).SetFloat(fl)
						default:
							//decimal
							dec, _ := strconv.ParseFloat(string(fl.([]byte)), 64)
							v.Field(i).SetFloat(dec)
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						switch in := item.(type) {
						case int64:
							v.Field(i).SetInt(in)
						default:
							dec, _ := strconv.ParseInt(string(in.([]byte)), 10, 64)
							v.Field(i).SetInt(dec)
						}
					case reflect.Struct:
						//time.Time
						switch tt := item.(type) {
						case time.Time:
							v.Field(i).Set(reflect.ValueOf(tt))
						default:
							tp, _ := time.Parse("15:04:05", string(tt.([]byte)))
							v.Field(i).Set(reflect.ValueOf(tp))
						}
					default:
						// SHOULD BE []byte
						v.Field(i).SetBytes(item.([]byte))
					}
				}
			}
		}
	}

	o = v.Interface()

	return
}
