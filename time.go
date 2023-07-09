package faketime

import (
	_ "unsafe"
)

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
	secondsPerWeek   = 7 * secondsPerDay
	daysPer400Years  = 365*400 + 97
	daysPer100Years  = 365*100 + 24
	daysPer4Years    = 365*4 + 1
)

//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

const (
	unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	wallToInternal int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
	internalToUnix int64 = -unixToInternal
)

var startNano int64 = runtimeNano() - 1

const (
	hasMonotonic = 1 << 63
	maxWall      = wallToInternal + (1<<33 - 1) // year 2157
	minWall      = wallToInternal               // year 1885
	nsecMask     = 1<<30 - 1
	nsecShift    = 30
)
