module github.com/AyushSenapati/guardian

go 1.15

replace github.com/AyushSenapati/limiter => ../limiter

require (
	github.com/AyushSenapati/limiter v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.1.2
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
)
