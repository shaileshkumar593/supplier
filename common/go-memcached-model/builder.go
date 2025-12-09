package gorm

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	ANDOperator = "AND"
	OROperator  = "OR"
)

var (
	JOIN_TYPES = map[string]bool{
		"LEFT":        true,
		"RIGHT":       true,
		"OUTER":       true,
		"INNER":       true,
		"LEFT OUTER":  true,
		"RIGHT OUTER": true,
	}

	QUERY = "SELECT %s FROM %s"
)

// Builder represents the core query builder
type Builder struct {
	DB             *sqlx.DB         `db:"-" json:"-"`
	Cache          *memcache.Client `db:"-" json:"-"`
	L              int              `db:"-" json:"-"`
	O              int              `db:"-" json:"-"`
	sortOrder      []string         `db:"-" json:"-"`
	CacheThis      bool             `db:"-" json:"-"`
	columns        []string         `db:"-" json:"-"`
	fromTable      string           `db:"-" json:"-"`
	groups         []string         `db:"-" json:"-"`
	joins          []JoinClause     `db:"-" json:"-"`
	wheres         []Where          `db:"-" json:"-"`
	bindings       []interface{}    `db:"-" json:"-"`
	query          string           `db:"-" json:"-"`
	bindVarCounter int32            `db:"-" json:"-"`
}

type Where struct {
	Type   string // OR / AND
	wheres []WhereClause
}

type WhereClause struct {
	Column   string      `db:"-" json:"-"`
	Operator string      `db:"-" json:"-"`
	Value    interface{} `db:"-" json:"-"`
	Type     string      `db:"-" json:"-"`
}

func (w *Where) Add(column string, op string, v interface{}, condOperator string) *Where {
	w.wheres = append(w.wheres, WhereClause{
		Column:   column,
		Operator: op,
		Value:    v,
		Type:     condOperator,
	})
	return w
}

type JoinClause struct {
	Condition string `db:"-" json:"-"`
	Table     string `db:"-" json:"-"`
	Type      string `db:"-" json:"-"`
}

// Connection can be used to change the default connection
func (b *Builder) Connection(c *sqlx.DB) *Builder {

	if c == nil || c.Ping() != nil {
		panic("Connection not available")
	}

	b.DB = c
	return b
}

// Columns to read form DB
func (b *Builder) Select(args ...string) *Builder {
	b.columns = args
	return b
}

// Set the left most table
func (b *Builder) Table(table string, alias string) *Builder {
	b.fromTable = table + " AS " + alias
	return b
}

// Join tables
func (b *Builder) Join(table string, tableAlias string, condition string, joinType string) *Builder {
	var join = JoinClause{}
	join.Table = table + " AS " + tableAlias
	join.Condition = condition
	table = strings.ToUpper(table)

	// check if in allowed joins
	joinType = strings.ToUpper(joinType)
	if _, ok := JOIN_TYPES[joinType]; ok {
		join.Type = joinType
	} else {
		panic("Unsupported join type")
	}

	b.joins = append(b.joins, join)

	return b
}

// Add a where clause with AND operator
func (b *Builder) Where(column string, op string, v interface{}) *Builder {
	b.whereClause(column, op, v, ANDOperator)
	return b
}

// Add a where clause with OR operator
func (b *Builder) ORWhere(column string, op string, v interface{}) *Builder {
	b.whereClause(column, op, v, OROperator)
	return b
}

// Where clause that is finally built and attached to the final query
func (b *Builder) whereClause(column string, op string, v interface{}, condOperator string) *Builder {
	var (
		w  = Where{}
		wc = WhereClause{}
	)
	w.Type = condOperator
	wc.Column = column
	wc.Operator = op
	wc.Value = v
	wc.Type = condOperator
	w.wheres = append(w.wheres, wc)
	b.wheres = append(b.wheres, w)

	return b
}

// Pass closure to build complex where clause with AND operator
//
//	Ex :WhereClosure(func() gorm.Where {
//				w := gorm.Where{}
//				w.Add("u.id", "=", 1, gorm.ANDOperator).
//					Add("u.email", "=", "user@email.com", gorm.OROperator)
//				return w
//			}).
//
// The closure should return a Where struct with multiple where clause in it
func (b *Builder) WhereClosure(c func() Where) *Builder {
	w := c()
	w.Type = ANDOperator
	b.wheres = append(b.wheres, w)
	return b
}

// Pass closure to build complex where clause with OR operator
//
//	Ex :ORWhereClosure(func() gorm.Where {
//				w := gorm.Where{}
//				w.Add("u.id", "=", 1, gorm.ANDOperator).
//					Add("u.email", "=", "user@email.com", gorm.OROperator)
//				return w
//			}).
//
// The closure should return a Where struct with multiple where clause in it
func (b *Builder) ORWhereClosure(c func() Where) *Builder {
	w := c()
	w.Type = OROperator
	b.wheres = append(b.wheres, w)
	return b
}

// Returns the final query with the bindings
func (b *Builder) ToSql() (sql string, bindings []interface{}) {
	if len(b.columns) > 0 {
		b.query = fmt.Sprintf(QUERY, strings.Join(b.columns, ", "), b.fromTable)
	} else {
		b.query = fmt.Sprintf(QUERY, "*", b.fromTable)
	}
	b.buildJoins()

	if len(b.wheres) > 0 {
		b.buildWhere()
	}

	return b.query, b.bindings
}

// Build joins from the join slice
func (b *Builder) buildJoins() {
	//TODO:improve joins
	// check if joins available
	if len(b.joins) > 0 {
		for _, j := range b.joins {
			b.query = b.query + " " + j.Type + " JOIN " + j.Table + " ON " + j.Condition
		}
	}
}

// Build where from the where slice
func (b *Builder) buildWhere() {
	b.query = b.query + " WHERE " // trailing space is necessary
	for _, w := range b.wheres {
		if b.query[len(b.query)-6:] == "WHERE " {
			w.Type = ""
		}
		// iterate over where clauses
		// if only one where clause then dont attach the brackets ()
		if len(w.wheres) > 1 {
			cw := " ("
			for _, wc := range w.wheres {
				b.bindings = append(b.bindings, wc.Value)
				if cw[len(cw)-1:] == "(" {
					wc.Type = ""
				}
				cw = cw + wc.Type + " " + wc.Column + " " + wc.Operator + " " + b.nextBindVar(wc.Value) + " "
			}
			b.query = b.query + w.Type + cw + ") "
		} else if len(w.wheres) == 1 {
			wc := w.wheres[0]
			wc.Type = w.Type // since its only one clause so the condition of the wrapper struct will be applied to it
			b.bindings = append(b.bindings, wc.Value)
			b.query = b.query + wc.Type + " " + wc.Column + " " + wc.Operator + " " + b.nextBindVar(wc.Value) + " "
		}
	}
}

// Returns the binding variable based on DB driver
func (b *Builder) nextBindVar(v interface{}) (bindvar string) {
	var (
		driver = b.DB.DriverName()
	)

	switch driver {
	case "mysql":
		bindvar = "?"
	case "postgres":
		b.bindVarCounter++
		bindvar = "$" + fmt.Sprint(b.bindVarCounter)
	}

	// check type and added quotes accordingly
	//t := reflect.TypeOf(v).String()
	//switch t {
	//case "string" :
	//	bindvar = `"` + bindvar + `"`
	//}
	return bindvar
}
