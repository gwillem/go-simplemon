package simplemon

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	maxLoad          = 0.75
	maxDiskPerc      = 0.9
	maxBackupAge     = 2 * 24 * time.Hour
	maxOpenFilesPerc = 0.8
	backupRoot       = "/backup"
)

func checkBackups() error {
	dirs := allDirsUnder(backupRoot)

	for _, dir := range dirs {
		pattern := filepath.Join(dir, "*")
		age, err := ageDaysOfNewestFile(pattern)
		if err != nil {
			return err
		}

		if age > maxBackupAge {
			return fmt.Errorf("newest file under %s is %.0f hours old", dir, age.Hours())
		}
	}
	return nil
}

func checkLoad() error {
	avg, err := load.Avg()
	if err != nil {
		return err
	}
	numCPU := runtime.NumCPU()
	if got := avg.Load5 / float64(numCPU); got > maxLoad {
		return fmt.Errorf("high load5 per cpu: %f", got)
	}
	return nil
}

func checkOpenFiles() error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		user, _ := p.Username()
		name, _ := p.Name()
		pname := fmt.Sprintf("%d/%s/%s", p.Pid, user, name)

		rlimits, err := p.Rlimit()
		if err != nil {
			continue
		}

		softLimit := rlimits[syscall.RLIMIT_NOFILE].Soft
		if softLimit <= 0 {
			// Skip processes with no file limits
			continue
		}

		cur, err := p.NumFDs()
		if err != nil || cur == 0 {
			continue
		}

		usage := float64(cur) / float64(softLimit)
		if usage > maxOpenFilesPerc {
			return fmt.Errorf("%s uses %d%% open files, are we growing too fast?", pname, int(usage*100))
		}
	}
	return nil
}

func checkDisk() error {
	parts, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	for _, part := range parts {
		if strings.Contains(part.Device, "loop") || strings.Contains(part.Mountpoint, "/snap/") ||
			strings.Contains(part.Mountpoint, "/boot") ||
			strings.Contains(part.Device, "devfs") {
			continue
		}

		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		// log.Printf("Disk %s bytes is %.0f%% full\n", part.Mountpoint, usage.UsedPercent)
		if usage.UsedPercent > 100*maxDiskPerc {
			return fmt.Errorf("disk %s bytes %.0f%% full", part.Mountpoint, usage.UsedPercent)
		}

		statvfs := syscall.Statfs_t{}
		err = syscall.Statfs(part.Mountpoint, &statvfs)
		if err != nil {
			continue
		}
		if statvfs.Files > 0 {
			percInodes := 100.0 * float64(statvfs.Files-statvfs.Ffree) / float64(statvfs.Files)
			// log.Printf("Disk %s inodes is %.0f%% full\n", part.Mountpoint, percInodes)
			if percInodes > 100*maxDiskPerc {
				return fmt.Errorf("disk %s inodes %.0f%% full", part.Mountpoint, percInodes)
			}
		}
	}
	return nil
}
