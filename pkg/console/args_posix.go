package console

import (
	"reflect"
	"strings"
	"unicode/utf8"
)

const posixArgNamesTag = "posix_args"
const posixOptionsTag = "posix_options"
const posixDefaultsTag = "posix_default"

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
	fields []posixField
	target reflect.Value
}

func newPosixFieldsParser(v interface{}) posixFieldsParser {
	parser := posixFieldsParser{
		target: reflect.ValueOf(v).Elem(),
	}
	t := reflect.TypeOf(v)
	var fields []posixField
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
			fields = append(fields, field)
		}
	}
	parser.fields = fields
	return parser
}

type posixArg struct {
	fieldName   string
	argType     string
	argNames    []string
	argDefault  string
	argRequired bool
	isSet       bool
}
type posixArgsParser struct {
	args   []posixArg
	target reflect.Value
}

func newPosixArgsParser(v interface{}) posixArgsParser {
	t := reflect.TypeOf(v)
	parser := posixArgsParser{target: reflect.ValueOf(v).Elem()}
	var args []posixArg
	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)
		if names, found := field.Tag.Lookup(posixArgNamesTag); found {
			arg := posixArg{
				fieldName:  field.Name,
				argType:    field.Type.Name(),
				argNames:   strings.Split(names, ","),
				argDefault: field.Tag.Get(posixDefaultsTag),
			}
			options := strings.Split(field.Tag.Get(posixOptionsTag), ",")
			for _, option := range options {
				if option == "required" {
					arg.argRequired = true
				}
			}
			args = append(args, arg)
		}
	}
	parser.args = args
	return parser
}
func (p posixArgsParser) findArg(arg string) []posixArg {
	var fields []posixArg
	for _, field := range p.args {
		for _, name := range field.argNames {
			if name == arg {
				fields = append(fields, field)
			}
		}
	}
	return fields
}
func (p posixArgsParser) Parse(args []string) error {
	var fields map[string]posixArg
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if utf8.RuneCountInString(arg) == 2 {
				continue
			}
			fields = make(map[string]posixArg)
			fList := p.findArg(strings.Replace(arg, "--", "", 1))
			for _, field := range fList {
				if field.argType == "bool" {
					p.target.FieldByName(field.fieldName).SetBool(true)
				}
				fields[field.fieldName] = field
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			argRunes := []rune(arg)
			if len(argRunes) == 1 {
				continue
			}
			fields = make(map[string]posixArg)
			for _, rune := range argRunes[1:] {
				fList := p.findArg(string(rune))
				for _, field := range fList {
					if field.argType == "bool" {
						p.target.FieldByName(field.fieldName).SetBool(true)
					}
					fields[field.fieldName] = field
				}
			}
		}
		for _, field := range fields {
			switch field.argType {
			case "string":
				p.target.FieldByName(field.fieldName).SetString("string")
			default:

			}
		}
	}
	return nil
}

func ParsePosixArgs(args []string, v interface{}) error {
	parser := newPosixArgsParser(v)
	return parser.Parse(args)
}
