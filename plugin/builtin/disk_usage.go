package builtin

import (
	"syscall"

	"github.com/jvikstedt/awake"
)

type DiskUsage struct{}

func (DiskUsage) Tag() string {
	return "builtin_disk_usage"
}

func (DiskUsage) Perform(scope awake.Scope) error {
	path, _ := scope.ValueAsString("path")

	disk, err := diskUsage(path)
	if err != nil {
		return err
	}

	scope.SetReturnVariable("all", awake.Variable{Type: awake.TypeFloat, Val: float64(disk.All) / float64(gb)})
	scope.SetReturnVariable("used", awake.Variable{Type: awake.TypeFloat, Val: float64(disk.Used) / float64(gb)})
	scope.SetReturnVariable("free", awake.Variable{Type: awake.TypeFloat, Val: float64(disk.Free) / float64(gb)})

	return nil
}

type diskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func diskUsage(path string) (diskStatus, error) {
	disk := diskStatus{}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return disk, err
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free

	return disk, nil
}

const (
	b  = 1
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
)
