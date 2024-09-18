package idgen

var Global IDGenerator = &KSUIDGenerator{}

type IDGenerator interface {
	Gen() string
}

// New returns a new unique ID.
func New[T ~string]() T {
	return T(Global.Gen())
}

func NewString() string {
	return New[string]()
}
