package statustype

type Type int

const (
	Readiness Type = 1 << iota
	Liveness
	Startup
)

func (target Type) Is(t Type) bool {
	return t&target == target
}

func Is(t Type, target Type) bool {
	return t&target == target
}
