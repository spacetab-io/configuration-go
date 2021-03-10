package stage

import (
	"fmt"
)

type Name string

const Defaults Name = "defaults"

func (sn Name) String() string {
	return fmt.Sprint(string(sn))
}
