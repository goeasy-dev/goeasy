package status

import (
	"context"
	"fmt"

	"goeasy.dev/status/statustype"
	"goeasy.dev/util"
)

const nameRandomLength = 8

func init() {
	checks = map[statustype.Type][]check{
		statustype.Readiness: make([]check, 0, 10),
		statustype.Liveness:  make([]check, 0, 10),
		statustype.Startup:   make([]check, 0, 10),
	}
}

type check struct {
	name string
	val  *bool
}

var checks map[statustype.Type][]check

func SimpleCheck(kind statustype.Type) *bool {
	caller := util.GetCaller()
	fmt.Println(caller)
	name := caller.Name
	if caller.Type != "" {
		name = fmt.Sprintf("%s.%s", caller.Type, caller.Name)
	}

	return NamedSimpleCheck(fmt.Sprintf("%s_%s", name, util.RandomString(nameRandomLength)), kind)
}

func NamedSimpleCheck(name string, t statustype.Type) *bool {
	v := false
	c := check{
		name: name,
		val:  &v,
	}

	if t.Is(statustype.Readiness) {
		checks[statustype.Readiness] = append(checks[statustype.Readiness], c)
	}
	if t.Is(statustype.Liveness) {
		checks[statustype.Liveness] = append(checks[statustype.Liveness], c)
	}
	if t.Is(statustype.Startup) {
		checks[statustype.Startup] = append(checks[statustype.Startup], c)
	}

	return c.val
}

func CheckStatus(ctx context.Context, kind statustype.Type) (map[string]bool, bool) {
	out := map[string]bool{}
	ok := true

	for _, check := range checks[kind] {
		out[check.name] = *check.val
		ok = ok && *check.val
	}

	return out, ok
}
