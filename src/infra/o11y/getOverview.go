package o11yInfra

import (
	"errors"
	"log"
	"os"
	"slices"
	"syscall"
	"time"

	"github.com/goinfinite/fleet/src/domain/entity"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	infraHelper "github.com/goinfinite/fleet/src/infra/helper"
	"github.com/shirou/gopsutil/disk"
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

func (repo GetOverview) getStorageUnitInfos() ([]valueObject.StorageUnitInfo, error) {
	var storageInfos []valueObject.StorageUnitInfo

	initialStats, err := disk.IOCounters()
	if err != nil {
		log.Printf("GetInitialStorageStatsFailed: %v", err)
		return storageInfos, errors.New("GetInitialStorageStatsFailed")
	}

	time.Sleep(time.Second)

	finalStats, err := disk.IOCounters()
	if err != nil {
		log.Printf("GetFinalStorageStatsFailed: %v", err)
		return storageInfos, errors.New("GetFinalStorageStatsFailed")
	}

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("GetPartitionsFailed: %v", err)
		return storageInfos, errors.New("GetPartitionsFailed")
	}

	desireableFileSystems := []string{
		"xfs",
		"btrfs",
		"ext4",
		"ext3",
		"ext2",
		"zfs",
		"vfat",
		"ntfs",
	}
	scannedDevices := []string{}
	for _, partition := range partitions {
		if !slices.Contains(desireableFileSystems, partition.Fstype) {
			continue
		}

		if slices.Contains(scannedDevices, partition.Device) {
			continue
		}

		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		initialStats := initialStats[partition.Device]
		finalStats := finalStats[partition.Device]

		deviceName, err := valueObject.NewDeviceName(partition.Device)
		if err != nil {
			continue
		}

		mountPoint, err := valueObject.NewUnixFilePath(partition.Mountpoint)
		if err != nil {
			continue
		}

		fileSystem, err := valueObject.NewUnixFileSystem(partition.Fstype)
		if err != nil {
			continue
		}

		readBytes := finalStats.ReadBytes - initialStats.ReadBytes
		readOpsCount := finalStats.ReadCount - initialStats.ReadCount
		writeBytes := finalStats.WriteBytes - initialStats.WriteBytes
		writeOpsCount := finalStats.WriteCount - initialStats.WriteCount

		storageUnitInfo := valueObject.NewStorageUnitInfo(
			deviceName,
			mountPoint,
			fileSystem,
			valueObject.Byte(usageStat.Total),
			valueObject.Byte(usageStat.Free),
			valueObject.Byte(usageStat.Used),
			usageStat.UsedPercent,
			valueObject.InodesCount(usageStat.InodesTotal),
			valueObject.InodesCount(usageStat.InodesFree),
			valueObject.InodesCount(usageStat.InodesUsed),
			usageStat.InodesUsedPercent,
			valueObject.Byte(readBytes),
			readOpsCount,
			valueObject.Byte(writeBytes),
			writeOpsCount,
		)

		storageInfos = append(storageInfos, storageUnitInfo)
		scannedDevices = append(scannedDevices, partition.Device)
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
	storageInfos    []valueObject.StorageUnitInfo
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
		storageInfos, err := repo.getStorageUnitInfos()
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
