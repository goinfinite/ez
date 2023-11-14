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
	storageInfos := []valueObject.StorageDeviceInfo{}
	for _, disk := range disks {
		if disk == "" {
			continue
		}

		deviceName, err := valueObject.NewDeviceName(disk)
		if err != nil {
			continue
		}

		storageInfo, err := repo.getStorageDeviceInfo(deviceName)
		if err != nil {
			return []valueObject.StorageDeviceInfo{}, errors.New("GetStorageDeviceInfoFailed")
		}

		storageInfos = append(storageInfos, storageInfo)
	}

	return storageInfos, nil
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

	cpuCores, err := valueObject.NewCpuCoresCount(len(cpuInfo))
	if err != nil {
		return hardwareSpecs, errors.New("GetCpuCoresCountFailed")
	}

	memoryLimit, err := repo.getMemLimit()
	if err != nil {
		return valueObject.HardwareSpecs{}, err
	}

	return valueObject.NewHardwareSpecs(
		cpuModel,
		cpuCores,
		cpuFrequency,
		valueObject.Byte(memoryLimit),
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

func (repo GetOverview) getNetInfos() ([]valueObject.NetInterfaceInfo, error) {
	var netInfos []valueObject.NetInterfaceInfo

	initialStats, err := net.IOCounters(true)
	if err != nil {
		log.Printf("GetInitialNetStatsFailed: %v", err)
		return netInfos, errors.New("GetInitialNetStatsFailed")
	}

	time.Sleep(time.Second)

	finalStats, err := net.IOCounters(true)
	if err != nil {
		log.Printf("GetFinalNetStatsFailed: %v", err)
		return netInfos, errors.New("GetFinalNetStatsFailed")
	}

	for i, interfaceStat := range finalStats {
		if interfaceStat.Name != initialStats[i].Name {
			continue
		}

		if interfaceStat.Name == "lo" {
			continue
		}

		deviceName, err := valueObject.NewDeviceName(interfaceStat.Name)
		if err != nil {
			continue
		}

		recvBytes := interfaceStat.BytesRecv - initialStats[i].BytesRecv
		recvPackets := interfaceStat.PacketsRecv - initialStats[i].PacketsRecv
		recvDropPackets := interfaceStat.Dropin - initialStats[i].Dropin
		recvErrs := interfaceStat.Errin - initialStats[i].Errin

		sentBytes := interfaceStat.BytesSent - initialStats[i].BytesSent
		sentPackets := interfaceStat.PacketsSent - initialStats[i].PacketsSent
		sentDropPackets := interfaceStat.Dropout - initialStats[i].Dropout
		sentErrs := interfaceStat.Errout - initialStats[i].Errout

		netInfo := valueObject.NewNetInterfaceInfo(
			deviceName,
			valueObject.Byte(recvBytes),
			recvPackets,
			recvDropPackets,
			recvErrs,
			valueObject.Byte(sentBytes),
			sentPackets,
			sentDropPackets,
			sentErrs,
		)

		netInfos = append(netInfos, netInfo)
	}

	return netInfos, nil
}

type HostResourceUsageResult struct {
	cpuUsagePercent float64
	memUsagePercent float64
	storageInfos    []valueObject.StorageDeviceInfo
	netInfos        []valueObject.NetInterfaceInfo
	err             error
}

func (repo GetOverview) getHostResourceUsage() (valueObject.HostResourceUsage, error) {
	cpuChan := make(chan HostResourceUsageResult)
	memChan := make(chan HostResourceUsageResult)
	storageChan := make(chan HostResourceUsageResult)
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
		storageInfos, err := repo.getStorageDeviceInfos()
		storageChan <- HostResourceUsageResult{storageInfos: storageInfos, err: err}
	}()

	go func() {
		netInfos, err := repo.getNetInfos()
		netChan <- HostResourceUsageResult{netInfos: netInfos, err: err}
	}()

	cpuResult := <-cpuChan
	if cpuResult.err != nil {
		return valueObject.HostResourceUsage{}, cpuResult.err
	}

	memResult := <-memChan
	if memResult.err != nil {
		return valueObject.HostResourceUsage{}, memResult.err
	}

	storageResult := <-storageChan
	if storageResult.err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetStorageInfoFailed")
	}

	netResult := <-netChan
	if netResult.err != nil {
		return valueObject.HostResourceUsage{}, errors.New("GetNetInfoFailed")
	}

	return valueObject.NewHostResourceUsage(
		cpuResult.cpuUsagePercent,
		memResult.memUsagePercent,
		storageResult.storageInfos,
		netResult.netInfos,
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
