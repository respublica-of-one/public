package posix

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const posixArgNamesTag = "posix_args"
const posixOptionsTag = "posix_options"
const posixDefaultsTag = "posix_default"

var ErrWrongFlagType = errors.New("wrong value type for flag")
var ErrNotImplementedType = errors.New("value setting for type not implemented")
var ErrCanNotParseValue = errors.New("can't parse value from string")
var ErrRequiredNotSet = errors.New("required field not set")
var ErrSettingDefaults = errors.New("can't set defaults")

func ParseArgs(v interface{}, args []string) error {
	parser := newPosixFieldsParser(v)
	if err := parser.setDefaults(); err != nil {
		return err
	}
	return parser.parse(args)
}

type posixOptions struct {
	touchDefaults  bool
	flag           bool
	required       bool
	defaultsString string
	names          []string
}
type posixField struct {
	fieldName  string
	fieldType  string
	options    posixOptions
	fieldIsSet bool
}
type posixFieldsParser struct {
	fields []*posixField
	target reflect.Value
}

func (p *posixFieldsParser) setValueFromString(field *posixField, value string) error {
	switch field.fieldType {
	case "int", "int8", "int16", "int32", "int64":
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			p.target.FieldByName(field.fieldName).SetInt(intVal)
		} else {
			return ErrCanNotParseValue
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
			p.target.FieldByName(field.fieldName).SetUint(uintVal)
		} else {
			return ErrCanNotParseValue
		}
	case "bool":
		if boolVal, err := strconv.ParseBool(value); err == nil {
			p.target.FieldByName(field.fieldName).SetBool(boolVal)
		} else {
			return ErrCanNotParseValue
		}
	case "float32":
		if floatVal, err := strconv.ParseFloat(value, 32); err == nil {
			p.target.SetFloat(floatVal)
		} else {
			return ErrCanNotParseValue
		}
	case "float64":
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			p.target.SetFloat(floatVal)
		} else {
			return ErrCanNotParseValue
		}
	case "string":
		p.target.FieldByName(field.fieldName).SetString(value)
	default:
		return ErrNotImplementedType
	}
	field.fieldIsSet = true
	return nil
}
func (p *posixFieldsParser) setDefaults() error {
	for _, field := range p.fields {
		if !field.options.touchDefaults && field.options.defaultsString != "" {
			if err := p.setValueFromString(field, field.options.defaultsString); err != nil {
				return ErrSettingDefaults
			}
		}
	}
	return nil
}
func (p *posixFieldsParser) checkRequired() error {
	for _, field := range p.fields {
		if field.options.required && !field.fieldIsSet {
			return ErrRequiredNotSet
		}
	}
	return nil
}
func (p *posixFieldsParser) parse(args []string) error {
	var fields []*posixField
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			fields = []*posixField{}
			argName := strings.TrimPrefix(arg, "--")
			field, err := p.touchField(argName)
			if err != nil {
				return err
			}
			if field != nil {
				fields = append(fields, field)
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			fields = []*posixField{}
			argNames := []rune(strings.TrimPrefix(arg, "-"))
			for _, argName := range argNames {
				field, err := p.touchField(string(argName))
				if err != nil {
					return err
				}
				if field != nil {
					fields = append(fields, field)
				}
			}
			continue
		}
		for _, field := range fields {
			if err := p.setValueFromString(field, arg); err != nil {
				return err
			}
		}
		fields = []*posixField{}
	}
	return p.checkRequired()
}
func (p *posixFieldsParser) touchField(arg string) (*posixField, error) {
	for _, field := range p.fields {
		for _, name := range field.options.names {
			if name == arg {
				if field.options.flag {
					return nil, p.touchFlag(field)
				}
				if field.options.touchDefaults && !field.fieldIsSet {
					if err := p.setValueFromString(field, field.options.defaultsString); err != nil {
						return nil, err
					}
					return field, nil
				}
				return field, nil
			}
		}
	}
	return nil, nil
}
func (p *posixFieldsParser) touchFlag(field *posixField) error {
	if field.fieldType != "bool" {
		return ErrWrongFlagType
	}
	p.target.FieldByName(field.fieldName).SetBool(true)
	return nil
}

func newPosixFieldsParser(v interface{}) *posixFieldsParser {
	parser := posixFieldsParser{
		target: reflect.ValueOf(v).Elem(),
	}
	t := reflect.TypeOf(v).Elem()
	var fields []*posixField
	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)
		if names, found := f.Tag.Lookup(posixArgNamesTag); found {
			field := posixField{
				fieldName: f.Name,
				fieldType: f.Type.String(),
			}
			field.options.names = strings.Split(names, ",")
			field.options.defaultsString = f.Tag.Get(posixDefaultsTag)
			options := strings.Split(f.Tag.Get(posixOptionsTag), ",")
			for _, option := range options {
				switch option {
				case "required":
					field.options.required = true
				case "touchDefaults":
					field.options.touchDefaults = true
				case "flag":
					field.options.flag = true
				}
			}
			fields = append(fields, &field)
		}
	}
	parser.fields = fields
	return &parser
}
