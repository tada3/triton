package tritondb

import "github.com/tada3/triton/logging"

var (
	log *logging.Entry
)

func init() {
	log = logging.NewEntry("tritondb")
}
