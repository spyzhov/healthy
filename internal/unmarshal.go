package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spyzhov/safe"
)

type InvalidType struct {
	Path  []string
	Value reflect.Value
	Given interface{}
}

func Unmarshal(path []string, target interface{}, value interface{}) (err error) {
	if safe.IsNil(value) {
		// do nothing...
		return nil
	}
	// region Setup
	ref := reflect.ValueOf(target)
	if ref.Kind() != reflect.Ptr || ref.IsNil() {
		panic(fmt.Sprintf("%s: wrong target type, ptr expected", strings.Join(path, ".")))
	}
	elem := ref.Elem()
	indi := elem
	if elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
		indi = elem.Elem()
	}
	if !indi.IsValid() {
		indi = elem
	}
	kind := indi.Kind()
	New := reflect.New(indi.Type())
	// endregion
	// region Unmarshaler
	unmarshaler, ok := New.Interface().(json.Unmarshaler)
	if ok {
		var data []byte
		if data, err = json.Marshal(value); err != nil {
			return fmt.Errorf("%s: can't set field: %w", strings.Join(path, "."), err)
		} else if err = unmarshaler.UnmarshalJSON(data); err != nil {
			return safe.Wrap(err, strings.Join(path, "."))
		}
		elem.Set(New.Elem())
		return nil
	}
	// endregion
	sValue, isStr := value.(string)
	// region Customizations for types
	switch kind {
	case reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer:
		return fmt.Errorf("%s: unsupported type `%s`", strings.Join(path, "."), kind)
	case reflect.Bool:
		if isStr {
			value, err = strconv.ParseBool(sValue)
			if err != nil {
				return safe.Wrap(err, strings.Join(path, "."))
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isStr {
			value, err = strconv.ParseInt(sValue, 10, 64)
			if err != nil {
				return safe.Wrap(err, strings.Join(path, "."))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if isStr {
			value, err = strconv.ParseUint(sValue, 10, 64)
			if err != nil {
				return safe.Wrap(err, strings.Join(path, "."))
			}
		}
	case reflect.Float32, reflect.Float64:
		if isStr {
			value, err = strconv.ParseFloat(sValue, 64)
			if err != nil {
				return safe.Wrap(err, strings.Join(path, "."))
			}
		}
	case reflect.Complex64, reflect.Complex128:
		// do nothing...
		break
	case reflect.Chan:
		// do nothing...
		break
	case reflect.Func:
		// do nothing...
		break
	case reflect.Map:
		val := reflect.ValueOf(getValue(value))
		if val.Kind() != reflect.Map {
			return &InvalidType{Path: path, Value: ref, Given: value}
		}
		eType := indi.Type().Elem()
		kType := indi.Type().Key()
		tMap := reflect.MakeMap(indi.Type())
		for _, key := range val.MapKeys() {
			kg := reflect.New(kType).Elem().Interface()
			tg := reflect.New(eType).Elem().Interface()
			iface := val.MapIndex(key).Interface()
			err = Unmarshal(append(path, fmt.Sprintf("key for `%v`", key.Interface())), &kg, key.Interface())
			if err != nil {
				return err
			}
			err = Unmarshal(append(path, fmt.Sprintf("[%v]", key.Interface())), &tg, iface)
			if err != nil {
				return err
			}
			tMap.SetMapIndex(reflect.ValueOf(kg), reflect.ValueOf(tg))
		}
		elem.Set(tMap)
		return nil
	case reflect.Array:
		val := reflect.ValueOf(getValue(value))
		if val.Kind() != reflect.Array || val.Len() != indi.Type().Len() {
			return &InvalidType{Path: path, Value: ref, Given: value}
		}
		tp := indi.Type().Elem()
		tArray := reflect.New(indi.Type()).Elem()
		for i := 0; i < val.Len(); i++ {
			tg := reflect.New(tp).Elem().Interface()
			iface := val.Index(i).Interface()
			err = Unmarshal(append(path, fmt.Sprintf("[%d]", i)), &tg, iface)
			if err != nil {
				return err
			}
			tArray.Index(i).Set(reflect.ValueOf(tg))
		}
		elem.Set(tArray)
		return nil
	case reflect.Slice:
		val := reflect.ValueOf(getValue(value))
		if val.Kind() != reflect.Slice {
			return &InvalidType{Path: path, Value: ref, Given: value}
		}
		tp := indi.Type().Elem()
		tSlice := reflect.MakeSlice(indi.Type(), val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			tg := reflect.New(tp).Elem().Interface()
			iface := val.Index(i).Interface()
			err = Unmarshal(append(path, fmt.Sprintf("[%d]", i)), &tg, iface)
			if err != nil {
				return err
			}
			tSlice.Index(i).Set(reflect.ValueOf(tg))
		}
		elem.Set(tSlice)
		return nil
	case reflect.Struct:
		val := reflect.ValueOf(getValue(value))
		if val.Kind() != reflect.Map {
			return &InvalidType{Path: path, Value: ref, Given: value}
		}
		kType := val.Type().Key()
		if kType.Kind() != reflect.String {
			return &InvalidType{Path: path, Value: ref, Given: value}
		}

		tp := indi.Type()
		tStruct := New.Elem()
		for i := 0; i < tp.NumField(); i++ {
			var iface interface{}
			var sPath []string
			field := tp.Field(i)
			if field.Anonymous {
				t := field.Type
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				if t.Kind() != reflect.Struct {
					continue
				}
				iface = value
				sPath = path
			} else {
				names := getNames(field)
				for _, name := range names {
					index := val.MapIndex(reflect.ValueOf(name))
					if index.IsValid() {
						iface = index.Interface()
						sPath = append(path, name)
						break
					}
				}
				if safe.IsNil(iface) { // default?
					continue
				}
			}
			fVal := reflect.New(field.Type).Elem().Interface()
			err = Unmarshal(sPath, &fVal, iface)
			if err != nil {
				return err
			}
			tStruct.Field(i).Set(reflect.ValueOf(fVal))
		}
		base(elem).Set(base(tStruct))
		return nil
	case reflect.Ptr:
		val := reflect.New(indi.Type().Elem())
		iface := val.Elem().Interface()
		err = Unmarshal(path, &iface, getValue(value))
		if err != nil {
			if _, ok := err.(*InvalidType); ok {
				return &InvalidType{Path: path, Value: ref, Given: value}
			}
			return err
		}
		val.Elem().Set(reflect.ValueOf(iface))
		elem.Set(val)
		return nil
	case reflect.Interface:
		elem.Set(reflect.ValueOf(value))
		return nil
	case reflect.String:
		// do nothing...
		break
	}
	// endregion
	// region Set
	rVal := reflect.ValueOf(value)
	if !rVal.Type().ConvertibleTo(elem.Elem().Type()) {
		return &InvalidType{Path: path, Value: ref, Given: value}
	}
	elem.Set(rVal.Convert(elem.Elem().Type()))
	// endregion
	return nil
}

func getNames(field reflect.StructField) []string {
	names := make([]string, 0, 4)
	name := field.Tag.Get("name")
	if name != "" {
		if name == "-" {
			return make([]string, 0)
		}
		names = append(names, name)
	}
	name = field.Tag.Get("json")
	if name != "" {
		parts := strings.Split(name, ",")
		if parts[0] == "-" {
			return make([]string, 0)
		}
		names = append(names, parts[0])
	}
	names = append(names, strcase.ToSnake(field.Name), field.Name)
	return names
}

func getValue(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	for {
		switch val.Kind() {
		case reflect.Interface, reflect.Ptr:
			val = val.Elem()
		default:
			return val.Interface()
		}
	}
}

func base(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func (t *InvalidType) Error() string {
	ref := reflect.Indirect(t.Value).Elem()
	tp := ref.Type()
	path := strings.Join(t.Path, ".")
	return fmt.Sprintf("%s: invalid type: expected `%s`, given `%T`", path, TypeToString(tp), t.Given)
}

func TypeToString(tp reflect.Type) string {
	types := ""
glob:
	for {
		switch tp.Kind() {
		case reflect.Ptr:
			types += "*"
			tp = tp.Elem()
		case reflect.Slice:
			types += "[]"
			tp = tp.Elem()
		case reflect.Array:
			types += fmt.Sprintf("[%d]", tp.Len())
			tp = tp.Elem()
		case reflect.Map:
			types += fmt.Sprintf("map[%s]", tp.Key().Kind())
			tp = tp.Elem()
		case reflect.Struct:
			types += tp.Name()
			break glob
		default:
			types += tp.Kind().String()
			break glob
		}
	}
	return types
}
