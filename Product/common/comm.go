package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		name := objValue.Type().Field(i).Name
		structFieldType := objValue.Field(i).Type()
		val := reflect.ValueOf(value)
		//fmt.Println(i, name, value, structFieldType, val)
		var err error
		if structFieldType != val.Type() {
			val, err = TypeConversion(value, structFieldType.Name())
			if err != nil {
			}
		}
		objValue.FieldByName(name).Set(val)
	}
}

func TypeConversion(value string, typeString string) (reflect.Value, error) {
	var v reflect.Value
	switch typeString {
	case "string":
		v = reflect.ValueOf(value)
	case "time.Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		v = reflect.ValueOf(t)
		if err != nil {
			return v, err
		}
	case "Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		v = reflect.ValueOf(t)
		if err != nil {
			return v, err
		}
	case "int":
		i, err := strconv.Atoi(value)
		v = reflect.ValueOf(i)
		if err != nil {
			return v, err
		}
	case "int8":
		i, err := strconv.ParseInt(value, 10, 64)
		v = reflect.ValueOf(int8(i))
		if err != nil {
			return v, err
		}
	case "int32":
		i, err := strconv.ParseInt(value, 10, 64)
		v = reflect.ValueOf(int32(i))
		if err != nil {
			return v, err
		}
	case "int64":
		i, err := strconv.ParseInt(value, 10, 64)
		v = reflect.ValueOf(int64(i))
		if err != nil {
			return v, err
		}
	case "float32":
		i, err := strconv.ParseFloat(value, 64)
		v = reflect.ValueOf(float32(i))
		if err != nil {
			return v, err
		}
	case "float64":
		i, err := strconv.ParseFloat(value, 64)
		v = reflect.ValueOf(i)
		if err != nil {
			return v, err
		}
		// add more
	default:
		v = reflect.ValueOf(value)
		return v, errors.New("Unsupported data type:" + typeString)
	}

	return v, nil
}
