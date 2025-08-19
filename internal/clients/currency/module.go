package currency

import "go.uber.org/fx"

var Module = fx.Module("currency",
	fx.Provide(
		CurrencyAdapter,
	),
)
