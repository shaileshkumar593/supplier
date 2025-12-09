package productConnections

import (
	"gopkg.in/guregu/null.v3"
)

// List represents an array of config on boarding model
type List []Model

// Model for formatting the query results
type Model struct {
	ID             null.String
	ProductCode    null.String
	DSN            null.String
	SlaveDSN       null.String
	ConnectionType null.String
	IsActive       null.String
}
