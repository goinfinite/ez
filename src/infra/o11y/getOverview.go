package o11yInfra

import (
	"errors"
	"log"
	"os"
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

func (repo GetOverview) getStorageDeviceInfo(
	deviceName valueObject.DeviceName,
) (valueObject.StorageDeviceInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return valueObject.StorageDeviceInfo{}, errors.New("StorageInfoError")
	}

	storageTotal := stat.Blocks * uint64(stat.Bsize)
	storageAvailable := stat.Bavail * uint64(stat.Bsize)
	storageUsed := storageTotal - storageAvailable

	return valueObject.NewStorageDeviceInfo(
		deviceName,
		valueObject.Byte(storageTotal),
		valueObject.Byte(storageAvailable),
		valueObject.Byte(storageUsed),
	), nil
}

func (repo GetOverview) getStorageDeviceInfos() ([]valueObject.StorageDeviceInfo, error) {
	disksList, err := infraHelper.RunCmd("lsblk", "-ndp", "-e", "7", "--output", "KNAME")
	if err != nil {
		log.Printf("GetDisksFailed: %v", err)
		return []valueObject.StorageDeviceInfo{}, errors.New("GetDisksFailed")
	}

	disks := strings.Split(disksList, "\n")
	deviceInfos := []valueObject.StorageDeviceInfo{}
	for _, disk := range disks {
		if disk == "" {
			continue
		}

		deviceName, err := valueObject.NewDeviceName(disk)
		if err != nil {
			continue
		}

		deviceInfo, err := repo.getStorageDeviceInfo(deviceName)
		if err != nil {
			return []valueObject.StorageDeviceInfo{}, errors.New("GetStorageDeviceInfoFailed")
		}

		deviceInfos = append(deviceInfos, deviceInfo)
	}

	return deviceInfos, nil
}

func (repo GetOverview) getMemLimit() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, errors.New("GetMemInfoFailed")
	}

	return memInfo.Total, nil
}

func (repo GetOverview) getHardwareSpecs() (valueObject.HardwareSpecs, error) {
	var hardwareSpecs valueObject.HardwareSpecs

	cpuInfo, err := cpu.Info()
	if err != nil {
		return hardwareSpecs, errors.New("GetCpuInfoFailed")
	}

	if len(cpuInfo) == 0 {
		return hardwareSpecs, errors.New("CpuInfoEmpty")
	}

	cpuModel, err := valueObject.NewCpuModelName(cpuInfo[0].ModelName)
	if err != nil {
		return hardwareSpecs, errors.New("GetCpuModelNameFailed")
	}

	cpuFrequency := cpuInfo[0].Mhz

	cpuCores, err := valueObject.NewCpuCoresCount(cpuInfo[0].Cores)
	if err != nil {
		return hardwareSpecs, errors.New("GetCpuCoresCountFailed")
	}

	memoryLimit, err := repo.getMemLimit()
	if err != nil {
		return valueObject.HardwareSpecs{}, err
	}

	storageInfo, err := repo.getStorageDeviceInfos()
	if err != nil {
		return valueObject.HardwareSpecs{}, err
	}

	return valueObject.NewHardwareSpecs(
		cpuModel,
		cpuCores,
		cpuFrequency,
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

type HostResourceUsageResult struct {
	cpuUsagePercent float64
	memUsagePercent float64
	deviceInfos     []valueObject.StorageDeviceInfo
	netUsage        valueObject.NetUsage
	err             error
}

func (repo GetOverview) getHostResourceUsage() (valueObject.HostResourceUsage, error) {
	cpuChan := make(chan HostResourceUsageResult)
	memChan := make(chan HostResourceUsageResult)
	deviceChan := make(chan HostResourceUsageResult)
	netChan := make(chan HostResourceUsageResult)

	go func() {
		cpuUsagePercent, err := repo.getCpuUsagePercent()
		cpuChan <- HostResourceUsageResult{cpuUsagePercent: cpuUsagePercent, err: err}
	}()

	go func() {
		memUsagePercent, err := repo.getMemUsagePercent()
		memChan <- HostResourceUsageResult{memUsagePercent: memUsagePercent, err: err}
	}()

	go func() {
		deviceInfos, err := repo.getStorageDeviceInfos()
		deviceChan <- HostResourceUsageResult{deviceInfos: deviceInfos, err: err}
	}()

	go func() {
		netUsage, err := repo.getNetUsage()
		netChan <- HostResourceUsageResult{netUsage: netUsage, err: err}
	}()

	cpuResult := <-cpuChan
	if cpuResult.err != nil {
		return valueObject.HostResourceUsage{}, cpuResult.err
	}

	memResult := <-memChan
	if memResult.err != nil {
		return valueObject.HostResourceUsage{}, memResult.err
	}

	deviceResult := <-deviceChan
	if deviceResult.err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetStorageInfoFailed")
	}

	netResult := <-netChan
	if netResult.err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetNetUsageFailed")
	}

	return valueObject.NewHostResourceUsage(
		cpuResult.cpuUsagePercent,
		memResult.memUsagePercent,
		deviceResult.deviceInfos,
		netResult.netUsage,
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
