package middleware

import (
	svc "swallow-supplier/iface"
)

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(svc.Service) svc.Service
