package console

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Args struct {
	fields   []string
	args     map[string]string
	defaults map[string]string
	ctx      map[string]string
	types    map[string]string
}

func NewArgs(val interface{}) Args {
	t := reflect.TypeOf(val)
	fieldsNumber := t.NumField()

	a := Args{
		fields:   make([]string, 0, fieldsNumber),
		args:     make(map[string]string),
		defaults: make(map[string]string),
		ctx:      make(map[string]string),
		types:    make(map[string]string),
	}

	for index := 0; index < fieldsNumber; index++ {
		field := t.Field(index)
		fieldName := field.Name
		a.fields = append(a.fields, fieldName)
		a.types[fieldName] = field.Type.Name()
		if args, found := field.Tag.Lookup("args"); found {
			argList := strings.Split(args, ",")
			for _, arg := range argList {
				a.args[arg] = fieldName
			}
		} else {
			a.args[fieldName] = fieldName
		}
		if defaultValue, found := field.Tag.Lookup("argsDefault"); found {
			a.defaults[fieldName] = defaultValue
		}
		if ctxKey, found := field.Tag.Lookup("ctx"); found {
			a.ctx[ctxKey] = fieldName
		}
	}
	return a
}

func (a Args) FromDefaults(val interface{}, skipFieldError bool) error {
	for key, value := range a.defaults {
		if err := a.setFieldFromString(val, key, value); err != nil && !skipFieldError {
			return err
		}
	}
	return nil
}

func (a Args) FromContext(val interface{}, ctx context.Context, skipFieldError bool) error {
	for key, field := range a.ctx {
		if value := ctx.Value(key); value != nil {
			if strVal, ok := value.(string); !ok && !skipFieldError {
				return errors.New(fmt.Sprintf("can't get ctx value string for %s : %+v", key, value))
			} else if strVal != "" {
				if err := a.setFieldFromString(val, field, strVal); err != nil && !skipFieldError {
					return err
				}
			}
		}
	}
	return nil
}

func (a Args) FromPosixArgs(val interface{}, args []string, skipFieldError bool) error {
	fmt.Printf("FromPosixArgs: %+v\n", args)
	var currentArgs []string
	valMode := false
	for _, arg := range args {
		fmt.Printf("vm: %t\targ: %s\t\n", valMode, arg)
		//full-command
		if strings.HasPrefix(arg, "--") {
			if valMode {
				if err := a.setFieldsFromString(val, currentArgs, "true"); err != nil && !skipFieldError {
					return err
				}
			}
			if fieldName := a.canonicalFieldName(strings.Replace(arg, "--", "", 1)); fieldName != "" {
				currentArgs = []string{fieldName}
				valMode = true
			}
			continue
		}
		//short-command list
		if strings.HasPrefix(arg, "-") {
			if valMode {
				if err := a.setFieldsFromString(val, currentArgs, "true"); err != nil && !skipFieldError {
					return err
				}
			}
			argsMappings := make(map[string]bool)
			for _, rune := range arg {
				if rune != '-' {
					if fieldName := a.canonicalFieldName(string(rune)); fieldName != "" {
						argsMappings[fieldName] = true
					}
				}
			}
			if len(argsMappings) > 0 {
				currentArgs = make([]string, 0, len(argsMappings))
				for key, _ := range argsMappings {
					currentArgs = append(currentArgs, key)
				}
				valMode = true
			}
			continue
		}
		//value
		if valMode {
			valMode = false
			fmt.Printf("%+v\n", currentArgs)
			if err := a.setFieldsFromString(val, currentArgs, arg); err != nil && !skipFieldError {
				return err
			}
		} else {
			if !skipFieldError {
				return errors.New("no field name available")
			}
		}
	}
	return nil
}

func (a Args) canonicalFieldName(source string) string {
	return a.args[source]
}

func (a Args) setFieldsFromString(val interface{}, fields []string, value string) error {
	var resErr error = nil
	for _, field := range fields {
		if err := a.setFieldFromString(val, field, value); err != nil {
			resErr = err
		}
	}
	return resErr
}

func (a Args) setFieldFromString(val interface{}, field, value string) error {
	elem := reflect.ValueOf(val).Elem()
	fmt.Printf("setFieldFromString: %s=%s\n", field, value)
	switch a.types[field] {
	case "int", "int8", "int16", "int32", "int64":
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		elem.FieldByName(field).SetInt(intVal)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		elem.FieldByName(field).SetUint(uintVal)
	case "bool":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		elem.FieldByName(field).SetBool(boolVal)
	default:
		elem.FieldByName(field).SetString(value)
	}
	return nil
}
