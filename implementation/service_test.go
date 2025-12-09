package implementation_test

import (
	"context"
	"net/http/httptest"

	"swallow-supplier/config"
	boilerplate "swallow-supplier/iface"
	service "swallow-supplier/iface"

	"github.com/go-kit/kit/log"
)

var (
	cf     config.AppConfig
	repo   map[string]boilerplate.MongoRepository
	ctx    = context.Background()
	srv    *httptest.Server
	logger log.Logger
	svc    service.Service
)

/*func TestMain(m *testing.M) {
	cf = *config.Instance()
	logger = log.NewNopLogger()
	_, repo = test.GetRepository(cf.DatabaseDriverPostgres, cf.DefaultDBName)
	srv = httptest.NewServer(test.GetHandlers())
	svc = implementation.NewService(repo, logger)
	defer srv.Close()
	os.Exit(m.Run())
}*/
