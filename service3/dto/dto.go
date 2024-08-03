package dto

import "time"

type ServiceResultDto struct {
	ExecutionTime   time.Time
	Service1Success bool
	Service2Success bool
	Service1Delay   int64
	Service2Delay   int64
}
