package middleware

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v3"

	customContext "swallow-supplier/context"
	customError "swallow-supplier/error"
	repo "swallow-supplier/iface"
	"swallow-supplier/utils"

	"swallow-supplier/config"
)

type (
	// Scope struct reprsents the scope, definition, and methods of a route
	Scope struct {
		ScopeList        []int            `json:"-"`
		Definition       string           `yaml:"def"`
		ProtectedMethods map[string][]int `yaml:"protected_methods,omitempty"`
	}

	// AuthConfig represents the route definitions
	AuthConfig struct {
		ScopeDefinition map[int]string   `yaml:"scope_definition"`
		RouteDefinition map[string]Scope `yaml:"route_definition"`
	}
)

var authConfig AuthConfig

// AuthMiddleWare for basic authentication
func AuthMiddleWare(repository map[string]repo.MongoRepository) endpoint.Middleware {
	// read the route definition
	authConfig, _ = ReadScopeDefinition()

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var logger = kitlog.NewJSONLogger(os.Stdout)
			c := config.Instance()

			logger = kitlog.With(logger,
				"ts", kitlog.DefaultTimestampUTC,
				"caller", kitlog.DefaultCaller,
			)

			// get the http method and route
			method := ctx.Value(httptransport.ContextKeyRequestMethod).(string)
			route := customContext.CtxRequestPath(ctx, customContext.CtxLabelRequestPathTemplate)
			level.Info(logger).Log(
				"type", "basic",
				"method", method,
				"route", route,
				customContext.CtxLabelRequestID, customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
			)

			// skip the default route
			if route == "" {
				return next(ctx, request)
			}

			// get route scope definition
			sc, err := GetRoute(route)
			if err != nil {
				err = customError.NewError(ctx, "route_not_defined", customError.ErrUndefinedRoute.Error(), nil)
				level.Error(logger).Log("error", err.Error())
				return nil, err
			}

			// check if the route method is protected and if not then just return
			if !sc.IsProtected(strings.ToUpper(method)) {
				return next(ctx, request)
			}

			ip, _ := GetClientIP(ctx)
			logger.Log(
				"client_ip", ip,
			)

			// get the context value of X-Api-Key header
			authToken := customContext.GetCtxHeader(ctx, customContext.CtxAuthorization)
			if authToken == "" {
				err = customError.NewErrorCustom(ctx, "401", "authorization key required", customError.ErrApiKeyRequired.Error(), http.StatusUnauthorized, "auth_middleware")
				level.Error(logger).Log("err", err)
				return nil, err
			}

			if c.AuthorizationKey != authToken {
				err = customError.NewErrorCustom(ctx, "401", "Authentication failed", customError.ErrApiKeyRequired.Error(), http.StatusUnauthorized, "auth_middleware")
				level.Error(logger).Log("err", err)
				return nil, err
			}

			return next(ctx, request)
		}
	}
}

// ReadScopeDefinition read the auth.yml
func ReadScopeDefinition() (conf AuthConfig, err error) {
	var (
		fc []byte
	)

	fc, err = os.ReadFile(config.Instance().AuthScopeConfigPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fc, &conf)
	if err != nil {
		return
	}

	return
}

// GetRoute - returns the auth route definition
func GetRoute(route string) (scope Scope, err error) {
	var (
		ok bool
	)

	if scope, ok = authConfig.RouteDefinition[route]; !ok {
		return scope, errors.New("invalid route")
	}

	return
}

// IsProtected - checks if the method is protected or not
func (s *Scope) IsProtected(method string) bool {
	var (
		ok bool
	)

	if s.ScopeList, ok = s.ProtectedMethods[method]; ok {
		return true
	}

	return false
}

// Verify - Call this method to verify if the consumer scope is in allowed scopes
func (s *Scope) Verify(consumerScope []int, method string) bool {
	// check if any of the scope matches with the route's allowed scope
	for _, uS := range consumerScope {
		for _, rS := range s.ScopeList {
			if uS == rS {
				return true
			}
		}
	}

	return false
}

// VerifyScope verify the consumer scope if allowed to access the resource
func VerifyScope(sc Scope, method string, scope string) string {
	scopeSplitted := strings.Split(scope, ",")
	scopeLevel := utils.SliceStringToSliceInt(scopeSplitted)

	if !sc.Verify(scopeLevel, strings.ToUpper(method)) {
		err := errors.New("access denied")
		return err.Error()
	}

	return ""
}

// GetClientIP returns the client ip
func GetClientIP(ctx context.Context) (string, error) {
	ipAddress := customContext.GetCtxHeader(ctx, customContext.CtxRemoteAddr)

	if ip := customContext.GetCtxHeader(ctx, customContext.CtxLabelForwardedFor); "" != ip {
		ipAddress = ip

		// X-Forwarded-For might contain multiple IPs. Get the last one.
		if strings.Contains(ipAddress, ",") {
			ips := strings.Split(ipAddress, ",")
			ipAddress = strings.Trim(ips[len(ips)-1], " ")
		}
	}

	var (
		ip  net.IP
		err error
	)

	if -1 != strings.Index(ipAddress, ":") {
		if ipAddress, _, err = net.SplitHostPort(ipAddress); nil != err {
			return "", err
		}
	}

	if err := ip.UnmarshalText([]byte(ipAddress)); nil != err {
		return "", err
	}

	return ipAddress, nil
}
