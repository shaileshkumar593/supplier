package gorm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	db "swallow-supplier/common/go-memcached-database"
	"swallow-supplier/common/go-resource/param"
	"swallow-supplier/common/go-tools/array"
	"swallow-supplier/common/go-tools/secure"
	"swallow-supplier/common/go-valid"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

const (
	// ISO8601Date ISO 8601 format with just the date
	ISO8601Date = "2006-01-02"

	// SQLDatetime YYYY-MM-DD HH:II:SS format with the date and time
	SQLDatetime = "2006-01-02 15:04:05"

	// TimestampFormat Timestamp Format
	TimestampFormat = "20060102150405"

	// MaxLimit ...
	MaxLimit = 50

	// MinOffset ...
	MinOffset = 0

	// SQLNullDate default null date
	SQLNullDate = "0000-00-00 00:00:00"

	// SQLNoDatabaseConnection No database connection
	SQLNoDatabaseConnection = "No database connection"

	// SQLInvalidOperator Invalid sql operator
	SQLInvalidOperator = "Invalid sql operator"

	//LimitRangeErrorCode limit must be a value 1 to 50
	LimitRangeErrorCode = "limit must be a value 1 to 50"

	//OffsetRangeErrorCode offset must be >= 0
	OffsetRangeErrorCode = "offset must be >= 0"

	// SQLNoRowsErrorCode sql: no rows in result set
	SQLNoRowsErrorCode = "sql: no rows in result set"

	// DateField ...
	DateField = "date"

	// JSONTag ...
	JSONTag = "json"

	// TypeString string
	TypeString = "string"

	// TypeInt int
	TypeInt = "int"

	// TypeFloat64 float64
	TypeFloat64 = "float64"
)

// OrderByIDDesc desc
var OrderByIDDesc = []string{"-id"}

// Model represents the core model
type Model struct {
	DB            *sqlx.DB         `db:"-" json:"-"`
	Txn           *sqlx.Tx         `db:"-" json:"-"`
	Cache         *memcache.Client `db:"-" json:"-"`
	Timezone      *string          `db:"-" json:"-"`
	Conditions    `db:"-" json:"-"`
	L             int      `db:"-" json:"-"`
	O             int      `db:"-" json:"-"`
	SortOrder     []string `db:"-" json:"-"`
	CacheThis     bool     `db:"-" json:"-"`
	IsTransaction bool     `db:"-" json:"-"`
	ServiceName   string
}

// UnixTimestamp return utc timestamp
func UnixTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Local().Unix())
}

// UnixToMysqlTime return utc timestamp
func UnixToMysqlTime(sec string, nsec string) string {
	iSec, _ := strconv.ParseInt(sec, 10, 64)
	iNsec, _ := strconv.ParseInt(nsec, 10, 64)

	return time.Unix(iSec, iNsec).Format(SQLDatetime)
}

// InitValidators Initializes the validations using github.com/go-validator/validator
func InitValidators() {
	validator.SetValidationFunc("pattern", valid.Pattern)
	validator.SetValidationFunc("sqlvalue", valid.SQLValue)
	validator.SetValidationFunc("required", valid.Required)
	validator.SetValidationFunc("greater", valid.GreaterThan)
	validator.SetValidationFunc("gender", valid.Gender)
	validator.SetValidationFunc("url", valid.URL)
}

// CacheGets performs a query using the sql and assigns to a list destination. caching can be applied
func (m *Model) CacheGets(dest interface{}, query string, args ...interface{}) error {
	var (
		err      error
		cacheKey string
	)

	if strings.Contains(query, "IN (?)") {
		query, args, err = sqlx.In(query, args...)

		if err != nil {
			return err
		}
	}

	query = m.DB.Rebind(query)
	query = m.ClearTicks(query)

	if m.CacheThis == true {
		argString, _ := json.Marshal(args)

		cacheKey = query + "_" + string(argString)
		if m.ServiceName != "" {
			cacheKey = m.ServiceName + "_" + cacheKey
		}

		if cache, _ := m.Cache.Get(secure.MD5(cacheKey)); cache != nil {
			json.Unmarshal(cache.Value, dest)
		} else {
			if err = m.DB.Select(dest, query, args...); err == nil {

				cacheKey = query + "_" + string(argString)
				if m.ServiceName != "" {
					cacheKey = m.ServiceName + "_" + cacheKey
				}

				cacheKey = secure.MD5(cacheKey)

				// Store child cache
				go func() {
					response, _ := json.Marshal(dest)
					m.Cache.Set(
						&memcache.Item{
							Key:        cacheKey,
							Value:      response,
							Expiration: int32(time.Now().Local().AddDate(1, 0, 0).Unix()),
						},
					)
				}()

				// Update parent cache
				go m.ParentCacheAddChildCache(query, cacheKey)
			}
		}
	} else {
		err = m.DB.Select(dest, query, args...)
	}

	return err
}

// CacheGet performs a query using the sql and assigns to destination. caching can be applied
func (m *Model) CacheGet(dest interface{}, query string) (err error) {
	var cacheKey string

	sql, arg := db.MapSQL(query, db.ToMap(dest))
	// Rebind bind vars depending upong database
	sql = m.DB.Rebind(sql)
	sql = m.ClearTicks(sql)

	if m.CacheThis == true {
		argString, _ := json.Marshal(arg)

		cacheKey = sql + "_" + string(argString)
		if m.ServiceName != "" {
			cacheKey = m.ServiceName + "_" + cacheKey
		}

		if cache, _ := m.Cache.Get(secure.MD5(cacheKey)); cache != nil {
			json.Unmarshal(cache.Value, dest)
		} else {
			if err = m.DB.Get(dest, sql, arg...); err == nil {
				cacheKey = sql + "_" + string(argString)
				if m.ServiceName != "" {
					cacheKey = m.ServiceName + "_" + cacheKey
				}

				cacheKey = secure.MD5(cacheKey)

				// Store child cache
				go func() {
					response, _ := json.Marshal(dest)
					m.Cache.Set(
						&memcache.Item{
							Key:        cacheKey,
							Value:      response,
							Expiration: int32(time.Now().Local().AddDate(1, 0, 0).Unix()),
						},
					)
				}()

				// Update parent cache
				go m.ParentCacheAddChildCache(sql, cacheKey)
			}
		}
	} else {
		err = m.DB.Get(dest, sql, arg...)
	}

	return
}

// ClearTicks removed ` (backtick) based on db driver
// currently it supports mysql and postgress
func (m *Model) ClearTicks(sql string) string {
	var driverName = m.DB.DriverName()

	switch driverName {
	case "postgres":
		return strings.Replace(sql, "`", "", -1)
	case "mysql":
		return sql
	default:
		return sql
	}
}

// CacheExec executes the sql.
// Caching not required cuz this method will only used to update data, not retrieve.
// Instead, in here we will flush caches associated to this table.
func (m *Model) CacheExec(sql string, arg interface{}) (sql.Result, error) {
	m.ParentCacheDeleteChildCache(sql)
	if m.IsTransaction {
		return m.Txn.NamedExec(sql, arg)
	}
	return m.DB.NamedExec(sql, arg)
}

// ConvertTimeStringToProductTime method - Parses Application datetime value to Product's timezone using "2006-01-02 15:04:05" format.
func (m *Model) ConvertTimeStringToProductTime(appDateTime string) string {
	tz, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation(SQLDatetime, appDateTime, tz)
	t, _ = m.ConvertToProductTime(t)

	return t.Format(SQLDatetime)
}

// ConvertTimeStringToAppTime method - Parses product datetime value to Application using "2006-01-02 15:04:05" format.
func (m *Model) ConvertTimeStringToAppTime(appDateTime string) string {
	tz, _ := time.LoadLocation(*m.Timezone)
	t, _ := time.ParseInLocation(SQLDatetime, appDateTime, tz)
	t, _ = m.ConvertToAppTime(t)

	return t.Format(SQLDatetime)
}

// ConvertToProductTime method - Converts passed Time object to current product's/program's time
func (m *Model) ConvertToProductTime(appDateTime time.Time) (productDateTime time.Time, err error) {
	tz, err := time.LoadLocation(*m.Timezone)
	productDateTime = appDateTime.In(tz)

	return
}

// ConvertToAppTime method - Converts passed Time object to Application time
func (m *Model) ConvertToAppTime(productDateTime time.Time) (appDateTime time.Time, err error) {
	tz, err := time.LoadLocation("Local")
	appDateTime = productDateTime.In(tz)

	return
}

// GenericUpdateSQL creates an sql update statement relative to the ERN
func (m *Model) GenericUpdateSQL(fields ...string) string {
	set := ""
	ctr := 1
	for _, v := range fields {
		set = set + "`" + v + "` = :" + v
		if ctr < len(fields) {
			set = set + ", "
		} else {
			set = set + " "
		}

		ctr++
	}

	return `UPDATE %s SET ` + set + "WHERE `ern` = :ern"
}

// GenericPrimaryKeyUpdateSQL creates an sql update statement relative to the Key
func (m *Model) GenericPrimaryKeyUpdateSQL(key string, fields ...string) string {
	set := ""
	ctr := 1
	for _, v := range fields {
		set += "`" + v + "` = :" + v
		if ctr < len(fields) {
			set += ", "
		} else {
			set += " "
		}

		ctr++
	}

	return `UPDATE %s SET ` + set + "WHERE `" + key + "` = :" + key
}

// JSON ...
func (m *Model) JSON(me interface{}) string {
	body, _ := json.Marshal(me)
	return string(body)
}

// Where set the conditions for the query
func (m *Model) Where(filters ...Condition) (err error) {
	m.Conditions = filters
	if _, _, err = m.Filters(filters, "?"); err != nil {
		return
	}

	return nil
}

// Limit set the limit for the query
func (m *Model) Limit(limit int) error {
	if limit < 1 || limit > 50 {
		return errors.New(LimitRangeErrorCode)
	}

	m.L = limit

	return nil
}

// Offset set the offset for the query
func (m *Model) Offset(offset int) error {

	if offset < 0 {
		return errors.New(OffsetRangeErrorCode)
	}

	m.O = offset

	return nil
}

// OrderBy set the sort order of the query
func (m *Model) OrderBy(order []string) {
	m.SortOrder = order
}

// PrepareStatement initialize the query
func (m *Model) PrepareStatement(model interface{}, fields ...string) (query string, args []interface{}) {
	var (
		filters                     []string
		order, filter, limit, field string
	)

	if len(fields) > 0 {
		field = fields[0]
	} else {
		field = ""
	}

	if len(m.SortOrder) > 0 {
		order = "ORDER BY " + strings.Join(m.SortOrder, ", ")
	}

	if len(m.Conditions) > 0 {
		m.Prepare(model, m.Conditions)
		filters, args, _ = m.Filters(m.Conditions, field)
		filter = "WHERE " + strings.Join(filters, " AND ")
	}

	if m.L != 0 {
		limit = fmt.Sprintf(`LIMIT %d OFFSET %d`, m.L, m.O)
	}

	query = filter + ` ` + order + ` ` + limit

	return
}

// GenericInsertSQL creates an sql insert statement
func (m *Model) GenericInsertSQL(fields ...string) string {
	set := ""
	values := ""
	ctr := 1
	for _, v := range fields {
		set = set + "`" + v + "`"
		values = values + ":" + v
		if ctr < len(fields) {
			set = set + ", "
			values = values + ", "
		} else {
			set = set + " "
			values = values + " "
		}

		ctr++
	}

	return `INSERT INTO %s (` + set + `) VALUES (` + values + `)`
}

// GenericDeleteSQL method - creates an sql delete statement relative to the ERN
func (m *Model) GenericDeleteSQL() string {
	return m.GenericDeleteRawSQL() + `WHERE ern = :ern`
}

// GenericPrimaryKeyDeleteSQL method - creates an sql delete statement relative to the Key
func (m *Model) GenericPrimaryKeyDeleteSQL(key string) string {
	return m.GenericDeleteRawSQL() + `WHERE ` + key + `= :` + key
}

// GenericDeleteRawSQL method - creates an sql delete statement
func (m *Model) GenericDeleteRawSQL() string {
	return `DELETE FROM %s `
}

// GetLastInsertedID simplifies the extraction of the LastInsertId from sql
func (m *Model) GetLastInsertedID(r sql.Result) string {
	if r == nil {
		return ""
	}

	id, _ := r.LastInsertId()

	if id == 0 {
		return ""
	}

	return strconv.FormatInt(id, 10)
}

// Filters ...
func (m *Model) Filters(filters Conditions, field string) (filter []string, args []interface{}, err error) {
	var (
		f, v, d string
		dates   []string
	)

	for _, cond := range filters {
		if err = cond.Validate(); err != nil {
			return
		}

		if "created_at" == cond.Field || "updated_at" == cond.Field {
			v = reflect.ValueOf(cond.Value).String()

			if dates, err = param.ParseDateRange(v); nil != err {
				return
			}

			f = fmt.Sprintf("`%s` BETWEEN ? AND ?", cond.Field)
			args = append(args, m.ConvertTimeStringToAppTime(dates[0]), m.ConvertTimeStringToAppTime(dates[1]))
			filter = append(filter, f)
		} else {
			if field == "" {
				d = fmt.Sprintf(":%s", cond.Field)
			} else {
				d = field
			}

			switch cond.Operator {
			case OperatorIn:
				f = fmt.Sprintf("`%s` %s (%s)", cond.Field, cond.Operator, d)
				args = append(args, cond.Value)

			case OperatorLike:
				v = reflect.ValueOf(cond.Value).String()
				v, err = param.ParseFilter(v, param.FilterRegexMatchAll)
				if nil != err {
					return
				}

				f = fmt.Sprintf("`%s` %s %s", cond.Field, cond.Operator, d)
				args = append(args, v)

			default:
				if reflect.TypeOf(cond.Value).String() == TypeInt {
					args = append(args, cond.Value)
				} else {
					args = append(args, reflect.ValueOf(cond.Value).String())
				}

				f = fmt.Sprintf("`%s` %s %s", cond.Field, cond.Operator, d)
			}

			filter = append(filter, f)
		}
	}

	return
}

// Prepare prepares model for query
func (m Model) Prepare(model interface{}, filters Conditions) {
	var (
		v    interface{}
		item string
	)

	for _, cond := range filters {
		item = cond.Field

		if reflect.TypeOf(cond.Value).String() == TypeInt {
			v = cond.Value
		} else {
			if OperatorBetween == strings.ToUpper(cond.Operator) {
				continue
			}

			if OperatorIn == strings.ToUpper(cond.Operator) {
				continue
			}

			v = reflect.ValueOf(cond.Value).String()

			if OperatorLike == strings.ToUpper(cond.Operator) {
				v, _ = param.ParseFilter(v.(string), param.FilterRegexWord)
			}
		}

		m.Assign(model, item, v)
	}
}

// Assign assigns the conditions to the model
func (m Model) Assign(model interface{}, search string, value interface{}, tags ...string) {
	var match string

	if len(tags) > 0 {
		match = tags[0]
	} else {
		match = "db"
	}

	t := reflect.ValueOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		items := strings.Split(f.Tag.Get(match), ",")

		if len(items) == 0 {
			continue
		}

		tag := items[0]

		if "" == tag || "-" == tag || DateField == tag {
			continue
		}

		if search == tag {
			if t.Type().Field(i).Type.String() == "sql.NullInt64" || t.Type().Field(i).Type.String() == "null.Int" {
				t.Field(i).FieldByName("Int64").SetInt(int64(value.(int)))
			} else if reflect.TypeOf(value).String() == TypeInt {
				t.Field(i).SetInt(int64(value.(int)))
			} else if t.Type().Field(i).Type.String() == "sql.NullFloat64" || t.Type().Field(i).Type.String() == "null.Float" {
				t.Field(i).FieldByName("Float64").SetFloat(value.(float64))
			} else if reflect.TypeOf(value).String() == TypeFloat64 {
				t.Field(i).SetFloat(value.(float64))
			} else if t.Type().Field(i).Type.String() == "sql.NullString" || t.Type().Field(i).Type.String() == "null.String" {
				t.Field(i).FieldByName("String").SetString(value.(string))
			} else {
				t.Field(i).SetString(value.(string))
			}
			return
		}
	}
}

// GetCondition create a confitions from all of the current field of the model
func (m Model) GetCondition(model interface{}) (conditions Conditions) {
	t := reflect.ValueOf(model).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		tag := f.Tag.Get("db")
		if "" == tag || "-" == tag {
			continue
		}

		if exist, _ := array.InArray(tag, []string{"id", "created_at", "updated_at"}); !exist && t.Field(i).String() != "" {
			conditions = append(conditions, Condition{
				Field:    tag,
				Operator: OperatorEqual,
				Value:    t.Field(i).String(),
			})
		}
	}

	return
}

// GenerateERN resource string DB Name, args[0] string programCode, args[1] string hash
func (m Model) GenerateERN(resource string, args ...string) (ern string, err error) {
	var (
		prefix, hash, service string
	)

	if os.Getenv("ERN_PREFIX") == "" {
		service = "fmt" // backward compatibility for FeeMS
	} else {
		service = os.Getenv("ERN_PREFIX")
	}

	if len(args) > 0 && args[0] != "" {
		programCode := args[0]
		prefix = fmt.Sprintf("%s:%s:%s", service, resource, programCode)
	} else {
		prefix = fmt.Sprintf("%s:%s", service, resource)
	}

	if len(args) > 1 && args[1] != "" {
		hash = args[1]
	} else {
		hash = fmt.Sprintf("%s", UnixTimestamp()+secure.RandomString(25, []rune(secure.RuneAlNumCS)))
	}

	ern = fmt.Sprintf("%s:%s", prefix, hash)

	return
}

// ParentCacheAddChildCache method - Adds new child cache to the master cache
func (m *Model) ParentCacheAddChildCache(sql string, childKey string) {
	key := getTableFromSQL(sql)
	m.ParentCacheSet(key, childKey)
}

// ParentCacheGet method - Gets its child keys
func (m *Model) ParentCacheGet(key string) string {
	if cache, _ := m.Cache.Get(key); cache != nil {
		return string(cache.Value)
	}

	return ""
}

// ParentCacheSet method - Sets new child keys
func (m *Model) ParentCacheSet(key string, newKey string) {
	cacheKeys := m.ParentCacheGet(key)

	if cacheKeys != "" {
		cacheKeys = cacheKeys + "," + newKey
	} else {
		cacheKeys = newKey
	}

	m.Cache.Set(
		&memcache.Item{
			Key:        key,
			Value:      []byte(cacheKeys),
			Expiration: int32(time.Now().Local().AddDate(1, 0, 0).Unix()),
		},
	)
}

// ParentCacheDeleteChildCache method - Delete master cache and its child caches
func (m *Model) ParentCacheDeleteChildCache(sql string) { // endpoint delete all child cache associated to the parent cache
	key := getTableFromSQL(sql)
	cacheKey := m.ParentCacheGet(key)

	for _, v := range strings.Split(cacheKey, ",") {
		m.Cache.Delete(v)
	}

	m.Cache.Delete(key)
}

// getTableFromSQL function - Returns table name from sql statement
func getTableFromSQL(sql string) (newSQL string) {
	words := strings.Split(sql, " ")

	for k, v := range words {
		if strings.ToUpper(v) == "FROM" {
			newSQL = words[k+1]
			break
		}

		if strings.ToUpper(v) == "INTO" {
			newSQL = words[k+1]
			break
		}

		if strings.ToUpper(v) == "UPDATE" {
			newSQL = words[k+1]
			break
		}
	}

	return newSQL
}
