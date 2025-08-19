package pgsql

import "go.uber.org/fx"

var Module = fx.Module("pgsql",
	fx.Provide(
		NewAuthRepo,
	),
)
