package http_test

import (
	"context"
	"net/http/httptest"

	"swallow-supplier/config"
	boilerplate "swallow-supplier/iface"
)

var (
	cf   config.AppConfig
	repo map[string]boilerplate.MongoRepository
	ctx  = context.Background()
	srv  *httptest.Server
)

/*func TestMain(m *testing.M) {
	cf = *config.Instance()
	_, repo = test.GetRepository( cf.DatabaseDriverPostgres, cf.DefaultDBName)
	srv = httptest.NewServer(test.GetHandlers())
	defer srv.Close()
	os.Exit(m.Run())
}*/
