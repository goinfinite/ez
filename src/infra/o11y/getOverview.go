package o11yInfra

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	infraHelper "github.com/goinfinite/fleet/src/infra/helper"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type GetOverview struct {
}

func (repo GetOverview) getUptime() (uint64, error) {
	sysinfo := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(sysinfo); err != nil {
		return 0, err
	}

	return uint64(sysinfo.Uptime), nil
}

func (repo GetOverview) getDiskInfo(
	diskName valueObject.DiskName,
) (valueObject.DiskInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return valueObject.DiskInfo{}, errors.New("StorageInfoError")
	}

	storageTotal := stat.Blocks * uint64(stat.Bsize)
	storageAvailable := stat.Bavail * uint64(stat.Bsize)
	storageUsed := storageTotal - storageAvailable

	return valueObject.NewDiskInfo(
		diskName,
		valueObject.Byte(storageTotal),
		valueObject.Byte(storageAvailable),
		valueObject.Byte(storageUsed),
	), nil
}

func (repo GetOverview) getDiskInfos() ([]valueObject.DiskInfo, error) {
	disksList, err := infraHelper.RunCmd("lsblk", "-ndp", "-e", "7", "--output", "KNAME")
	if err != nil {
		log.Printf("GetDisksFailed: %v", err)
		return []valueObject.DiskInfo{}, errors.New("GetDisksFailed")
	}

	disks := strings.Split(disksList, "\n")
	diskInfos := []valueObject.DiskInfo{}
	for _, disk := range disks {
		if disk == "" {
			continue
		}

		diskName, err := valueObject.NewDiskName(disk)
		if err != nil {
			continue
		}

		diskInfo, err := repo.getDiskInfo(diskName)
		if err != nil {
			return []valueObject.DiskInfo{}, errors.New("GetDiskInfoFailed")
		}

		diskInfos = append(diskInfos, diskInfo)
	}

	return diskInfos, nil
}

func (repo GetOverview) getMemLimit() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, errors.New("GetMemInfoFailed")
	}

	return memInfo.Total, nil
}

func (repo GetOverview) getHardwareSpecs() (valueObject.HardwareSpecs, error) {
	cmd := exec.Command(
		"awk",
		"-F:",
		"/vendor_id/{vendor=$2} /cpu MHz/{freq=$2} END{print vendor freq}",
		"/proc/cpuinfo",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("GetCpuSpecsFailed: %v", err)
		return valueObject.HardwareSpecs{}, errors.New("GetCpuSpecsFailed")
	}
	trimmedOutput := strings.TrimSpace(string(output))
	if trimmedOutput == "" {
		return valueObject.HardwareSpecs{}, errors.New("EmptyCpuSpecs")
	}

	cpuInfo := strings.Split(trimmedOutput, " ")
	if len(cpuInfo) < 2 {
		return valueObject.HardwareSpecs{}, errors.New("ParseCpuSpecsFailed")
	}

	cpuModel := strings.TrimSpace(cpuInfo[0])
	cpuFrequency := strings.TrimSpace(cpuInfo[1])
	cpuFrequencyFloat, err := strconv.ParseFloat(cpuFrequency, 64)
	if err != nil {
		log.Printf("GetCpuFrequencyFailed: %v", err)
		return valueObject.HardwareSpecs{}, errors.New("GetCpuFrequencyFailed")
	}

	cpuCoresStr, err := infraHelper.RunCmd("nproc")
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetCpuCoresStrFailed")
	}
	cpuCores, err := strconv.ParseUint(cpuCoresStr, 10, 64)
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetCpuCoresFailed")
	}

	memoryLimit, err := repo.getMemLimit()
	if err != nil {
		return valueObject.HardwareSpecs{}, err
	}

	storageInfo, err := repo.getDiskInfos()
	if err != nil {
		return valueObject.HardwareSpecs{}, err
	}

	return valueObject.NewHardwareSpecs(
		cpuModel,
		cpuCores,
		cpuFrequencyFloat,
		valueObject.Byte(memoryLimit),
		storageInfo,
	), nil
}

func (repo GetOverview) getCpuUsagePercent() (float64, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, errors.New("GetCpuUsageFailed")
	}

	return cpuPercent[0], nil
}

func (repo GetOverview) getMemUsagePercent() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, errors.New("GetMemInfoFailed")
	}

	return memInfo.UsedPercent, nil
}

func (repo GetOverview) getNetUsage() (valueObject.NetUsage, error) {
	initialStats, err := net.IOCounters(true)
	if err != nil {
		log.Printf("GetInitialNetStatsFailed: %v", err)
		return valueObject.NetUsage{}, errors.New("GetInitialNetStatsFailed")
	}

	time.Sleep(time.Second)

	finalStats, err := net.IOCounters(true)
	if err != nil {
		log.Printf("GetFinalNetStatsFailed: %v", err)
		return valueObject.NetUsage{}, errors.New("GetFinalNetStatsFailed")
	}

	sentBytes := finalStats[0].BytesSent - initialStats[0].BytesSent
	recvBytes := finalStats[0].BytesRecv - initialStats[0].BytesRecv

	return valueObject.NewNetUsage(
		valueObject.Byte(sentBytes),
		valueObject.Byte(recvBytes),
	), nil
}

func (repo GetOverview) getHostResourceUsage() (
	valueObject.HostResourceUsage,
	error,
) {
	cpuUsagePercent, err := repo.getCpuUsagePercent()
	if err != nil {
		return valueObject.HostResourceUsage{}, err
	}
	memUsagePercent, err := repo.getMemUsagePercent()
	if err != nil {
		return valueObject.HostResourceUsage{}, err
	}

	diskInfos, err := repo.getDiskInfos()
	if err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetStorageInfoFailed")
	}

	netUsage, err := repo.getNetUsage()
	if err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetNetUsageFailed")
	}

	return valueObject.NewHostResourceUsage(
		cpuUsagePercent,
		memUsagePercent,
		diskInfos,
		netUsage,
	), nil
}

func (repo GetOverview) Get() (entity.O11yOverview, error) {
	hostnameStr, err := os.Hostname()
	if err != nil {
		hostnameStr = "localhost"
	}

	hostname, err := valueObject.NewFqdn(hostnameStr)
	if err != nil {
		return entity.O11yOverview{}, errors.New("GetHostnameFailed")
	}

	uptime, err := repo.getUptime()
	if err != nil {
		uptime = 0
	}

	publicIpAddress, err := infraHelper.GetPublicIpAddress()
	if err != nil {
		publicIpAddress, _ = valueObject.NewIpAddress("0.0.0.0")
	}

	hardwareSpecs, err := repo.getHardwareSpecs()
	if err != nil {
		log.Printf("GetHardwareSpecsFailed: %v", err)
		return entity.O11yOverview{}, errors.New("GetHardwareSpecsFailed")
	}

	currentResourceUsage, err := repo.getHostResourceUsage()
	if err != nil {
		log.Printf("GetHostResourceUsageFailed: %v", err)
		return entity.O11yOverview{}, errors.New("GetHostResourceUsageFailed")
	}

	return entity.NewO11yOverview(
		hostname,
		uptime,
		publicIpAddress,
		hardwareSpecs,
		currentResourceUsage,
	), nil
}
