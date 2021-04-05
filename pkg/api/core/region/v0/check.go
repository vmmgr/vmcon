package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
)

func check(input core.Region) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	return nil
}
