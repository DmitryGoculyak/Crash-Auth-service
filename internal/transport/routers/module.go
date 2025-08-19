package routers

import "go.uber.org/fx"

var Module = fx.Module("routers",
	fx.Invoke(
		RunServer,
	),
)
