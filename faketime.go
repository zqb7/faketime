package faketime

import (
	"sync"
	"time"

	"bou.ke/monkey"
)

var (
	patchGuardTimeNow *monkey.PatchGuard
)

type fakeTime struct {
	sync.Once
	wall uint64
	ext  int64

	fakeSub int64
}

func (f *fakeTime) sec() int64 {
	if f.wall&hasMonotonic != 0 {
		return wallToInternal + int64(f.wall<<1>>(nsecShift+1))
	}
	return f.ext
}

func (f *fakeTime) unixSec() int64 { return f.sec() + internalToUnix }

func (f *fakeTime) nsec() int32 { return int32(f.wall & nsecMask) }

// 微秒
func (f *fakeTime) unixMicro() int64 { return f.unixSec()*1e6 + int64(f.nsec())/1e3 }

// 毫秒
func (f *fakeTime) unixMilli() int64 { return f.unixSec()*1e3 + int64(f.nsec())/1e6 }

// 保持当前时间
func (f *fakeTime) realTime() time.Time {
	sec, nsec, mono := now()
	mono -= startNano
	sec += unixToInternal - wallToInternal
	wall := hasMonotonic | uint64(sec)<<nsecShift | uint64(nsec)
	if uint64(sec)>>33 != 0 {
		wall = uint64(nsec)
	}
	f.wall = wall
	return time.UnixMilli(f.unixMilli())
}

func (f *fakeTime) Add(d time.Duration) time.Time {
	return f.realTime().Add(d)
}

// 从固定的一个时间开始
func (f *fakeTime) FixTime(year, month, day, hour, min, sec int) time.Time {
	f.Once.Do(func() {
		f.fakeSub = time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local).UnixMilli() - f.realTime().UnixMilli()
	})
	return f.realTime().Add(time.Millisecond * time.Duration(f.fakeSub))
}

func RealTime() {
	ftime := &fakeTime{}
	f := func() time.Time { return ftime.realTime() }
	patchTimeNow(time.Now, func() time.Time {
		return f()
	})
}

func Add(d time.Duration) {
	ftime := &fakeTime{}
	f := func() time.Time { return ftime.Add(d) }
	patchTimeNow(time.Now, func() time.Time {
		return f()
	})
}

func FixTime(year, month, day, hour, min, sec int) {
	ftime := &fakeTime{}
	f := func() time.Time { return ftime.FixTime(year, month, day, hour, min, sec) }
	patchTimeNow(time.Now, func() time.Time {
		return f()
	})
}

func patchTimeNow(target, replacement interface{}) {
	if patchGuardTimeNow != nil {
		patchGuardTimeNow.Unpatch()
	}
	patchGuardTimeNow = Patch(target, replacement)
}

func Patch(target, replacement interface{}) *monkey.PatchGuard {
	return monkey.Patch(target, replacement)
}

func UnpatchAll() {
	monkey.UnpatchAll()
}
