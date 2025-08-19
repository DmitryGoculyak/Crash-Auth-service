package billing

import "go.uber.org/fx"

var Module = fx.Module("billing",
	fx.Provide(
		BillingAdapter,
	),
)
