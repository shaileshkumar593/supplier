package db

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	typeSQLNullBool   = "sql.NullBool"
	typeSQLNullFloat  = "sql.NullFloat64"
	typeSQLNullInt    = "sql.NullInt64"
	typeSQLNullString = "sql.NullString"
	typeNullBool      = "null.Bool"
	typeNullFloat     = "null.Float"
	typeNullInt       = "null.Int"
	typeNullString    = "null.String"
)

// MapSQL maps the :token to the model field value
func MapSQL(query string, model map[string]interface{}) (string, []interface{}) {
	var val []interface{}
	re := regexp.MustCompile(`(^|\s+):\w+([\s]+|$)`)

	sql := re.ReplaceAllStringFunc(query,
		func(match string) string {
			k := strings.Trim(match, ": \r\n\t")
			if v, ok := model[k]; ok {
				val = append(val, v)
				return strings.Replace(match, ":"+k, "?", -1)
			}

			return match
		})

	return sql, val
}

// ToMap returns the value of the the fields with db tag
func ToMap(model interface{}) map[string]interface{} {
	var (
		f        reflect.StructField
		tag      string
		modelMap map[string]interface{}
	)

	modelMap = make(map[string]interface{})
	t := reflect.ValueOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		f = t.Type().Field(i)
		tag = f.Tag.Get("db")

		if "" == tag || "-" == tag {
			continue
		}

		switch t.Field(i).Kind().String() {
		case "int":
			modelMap[tag] = t.Field(i).Int()
		case "struct":
			switch t.Field(i).Type().String() {
			case typeSQLNullBool:
				modelMap[tag] = t.Field(i).Field(0).Bool()
			case typeSQLNullFloat:
				modelMap[tag] = t.Field(i).Field(0).Float()
			case typeSQLNullInt:
				modelMap[tag] = t.Field(i).Field(0).Int()
			case typeSQLNullString:
				modelMap[tag] = t.Field(i).Field(0).String()
			case typeNullBool:
				modelMap[tag] = t.Field(i).Field(0).Field(0).Bool()
			case typeNullFloat:
				modelMap[tag] = t.Field(i).Field(0).Field(0).Float()
			case typeNullInt:
				modelMap[tag] = t.Field(i).Field(0).Field(0).Int()
			case typeNullString:
				modelMap[tag] = t.Field(i).Field(0).Field(0).String()
			default:
				modelMap[tag] = t.Field(i).String()
			}
		default:
			modelMap[tag] = t.Field(i).String()
		}
	}

	return modelMap
}
