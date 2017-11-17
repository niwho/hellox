package strategy

import (
	"github.com/niwho/hellox/logs"
)

type IS interface {
	GetIntoMatched(uid string) string
}
