package main

import (
	"syscall"
)

// cpuCycles returns the number of CPU cycles
func cpuCycles() int64 {
	usage := new(syscall.Rusage)
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, usage)
	return usage.Utime.Nano() + usage.Stime.Nano()
}
