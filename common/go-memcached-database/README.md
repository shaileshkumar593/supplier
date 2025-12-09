# db
--
    import "bitbucket.org/matchmove/go-memcached-database"


## Usage

#### func  MapSQL

```go
func MapSQL(query string, model map[string]string) (string, []interface{})
```
MapSQL maps the :token to the model field value

#### func  ToMap

```go
func ToMap(model interface{}) map[string]string
```
ToMap returns the value of the the fields with db tag

#### type DB

```go
type DB struct {
	Connection *sqlx.DB
	Status     string
	Driver     string
	Open       string
}
```

DB represents database config

#### func  New

```go
func New(driver string, open string) (*DB, error)
```
New creates a new (DB)Database object

#### func (*DB) Close

```go
func (d *DB) Close()
```
Close sets up a connection using the current credentials

#### func (*DB) Connect

```go
func (d *DB) Connect() (*sqlx.DB, error)
```
Connect sets up a connection using the current credentials

#### type Databases

```go
type Databases map[string]*DB
```

Databases maps all databases used by this application
