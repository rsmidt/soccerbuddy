package eventing

type VersionMatcher interface {
	Matches(current, actual AggregateVersion) bool
}

type VersionMatcherFunc func(current, actual AggregateVersion) bool

func (f VersionMatcherFunc) Matches(current, actual AggregateVersion) bool {
	return f(current, actual)
}

var (
	VersionMatcherAlways VersionMatcher = versionMatcherAlways{}
	VersionMatcherExact  VersionMatcher = VersionMatcherFunc(versionMatcherExact)
)

type versionMatcherAlways struct{}

func (versionMatcherAlways) Matches(_, _ AggregateVersion) bool {
	return true
}

func versionMatcherExact(current, actual AggregateVersion) bool {
	return current == actual
}
