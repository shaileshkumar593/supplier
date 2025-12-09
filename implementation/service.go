package implementation

import (
	"swallow-supplier/config"
	svc "swallow-supplier/iface"

	"github.com/go-kit/kit/log"
)

var (
	cf = config.Instance()
)

// service implements the Order Service
type service struct {
	mongoRepository map[string]svc.MongoRepository
	logger          log.Logger
}

// NewService creates and returns a new service instance
func NewService(mongorepo map[string]svc.MongoRepository, logger log.Logger) svc.Service {
	return &service{
		mongoRepository: mongorepo,
		logger:          logger,
	}
}

// FormatCompiledErrors formats error response
func (s *service) FormatCompiledErrors(errs []interface{}) map[string][]interface{} {
	var errSource = make(map[string][]interface{})
	for _, e := range errs {
		if e.(map[string][]interface{}) != nil {
			for f, ex := range e.(map[string][]interface{}) {
				errSource[f] = ex
			}
		}
	}

	return errSource
}
