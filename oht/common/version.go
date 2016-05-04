package common

import (
	"fmt"
)

func CompileClientInfo(name, version string) string {
	return fmt.Sprintf("%s/v%s", name, version)
}
