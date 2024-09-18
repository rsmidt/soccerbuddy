package eventing

type QueryConfig struct {
	ErrorOnCryptoShredded           bool
	LimitToOldestRunningTransaction bool
}

type QueryOpts interface {
	Apply(config QueryConfig) QueryConfig
}

type queryOptsFunc func(config QueryConfig) QueryConfig

func (f queryOptsFunc) Apply(config QueryConfig) QueryConfig {
	return f(config)
}

func WithErrorOnCryptoShredded() QueryOpts {
	return queryOptsFunc(func(config QueryConfig) QueryConfig {
		config.ErrorOnCryptoShredded = true
		return config
	})
}

// WithLimitToOldestRunningTransaction instructs the query to only return rows that were inserted
// before the oldest running transaction. This ensures that, e.g., projectors do not skip events.
func WithLimitToOldestRunningTransaction() QueryOpts {
	return queryOptsFunc(func(config QueryConfig) QueryConfig {
		config.LimitToOldestRunningTransaction = true
		return config
	})
}
