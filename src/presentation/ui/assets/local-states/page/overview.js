document.addEventListener('alpine:init', () => {
    Alpine.data('containers', () => ({
        // Primary State
        container: {},
        resetPrimaryState() {
            this.container = {
                'id': '',
                'accountId': '',
                'hostname': '',
                'status': true,
                'imageAddress': 'docker.io/goinfinite/os:latest',
                'portBindings': [],
                'restartPolicy': 'unless-stopped',
                'entrypoint': '',
                'profileId': "1",
                'envs': [],
                'launchScript': '',
                'autoCreateMappings': true,
                'existingContainerId': '',
                'useImageExposedPorts': false,
            }
        },
        init() {
            this.resetPrimaryState()
        },
        async autoPortBindingsRetriever(prevImageAddress, newImageAddress) {
            if (!this.isCreateContainerModalOpen) {
                return
            }

            const isSameImage = prevImageAddress === newImageAddress
            const isPortBindingsEmpty = this.container.portBindings.length === 0
            if (isSameImage && !isPortBindingsEmpty) {
                return
            }

            const isImageAddressTooShort = newImageAddress.length < 3
            const isImageNotAddressYet = !newImageAddress.includes('/')
            if (isImageAddressTooShort || isImageNotAddressYet) {
                this.container.portBindings = []
                return
            }

            const remoteApiUrl = '/api/v1/container/registry/image/tagged/?address=' + this.container.imageAddress
            this.container.portBindings = await fetch(remoteApiUrl, {
                method: "GET",
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
            })
                .then((apiResponse) => {
                    if (!apiResponse.ok) {
                        throw new Error('BadHttpResponseCode: ' + apiResponse.status)
                    }

                    return apiResponse.json()
                })
                .then((jsonResponse) => {
                    return jsonResponse.body.portBindings
                })
                .catch((error) => {
                    console.error('AutoUpdatePortBindingsError: ', error)
                    return []
                })
        },
        syncExistingContainerPortBindings(prevId, newId) {
            if (!this.isCreateContainerModalOpen) {
                return
            }

            const isSameContainer = prevId === newId
            const isPortBindingsEmpty = this.container.portBindings.length === 0
            if (isSameContainer && !isPortBindingsEmpty) {
                return
            }

            if (newId.length < 12) {
                this.container.portBindings = []
                return
            }

            const containerEntity = JSON.parse(
                document.getElementById('containerEntity_' + newId).textContent
            )
            portBindingsWithoutPrivatePorts = containerEntity.portBindings.map(
                originalPortBinding => {
                    adjustedPortBinding = Object.assign({}, originalPortBinding)

                    adjustedPortBinding.publicPort = adjustedPortBinding.containerPort
                    delete adjustedPortBinding.privatePort

                    return adjustedPortBinding
                },
            )

            this.container.portBindings = portBindingsWithoutPrivatePorts
        },

        // Auxiliary States
        createContainerSourceSelectedOption: 'Marketplace',
        createContainerSelectedMarketplaceType: 'App',
        createContainerSelectedMarketplaceItemId: 0,
        updateSelectedMarketplaceItem(itemId) {
            this.createContainerSelectedMarketplaceItemId = itemId
            const marketplaceItemEntity = JSON.parse(
                document.getElementById('marketplaceItemEntity_' + itemId).textContent
            )
            this.container.imageAddress = marketplaceItemEntity.registryImageAddress
            this.container.launchScript = ''
            if (marketplaceItemEntity.hasOwnProperty('launchScript')) {
                this.container.launchScript = marketplaceItemEntity.launchScript
            }
        },
        containersFilters: {
            containerId: '',
            accountId: '',
            hostname: '',
            status: '',
            imageId: '',
            imageAddress: '',
            imageHash: '',
            restartPolicy: '',
            profileId: '',
        },
        containersPagination: {
            pageNumber: containersCurrentPageNumber,
            itemsPerPage: 5,
        },
        reloadContainersTable() {
            queryParams = new URLSearchParams()
            queryParams.set('containersPageNumber', this.containersPagination.pageNumber)
            queryParams.set('containersItemsPerPage', this.containersPagination.itemsPerPage)

            for (let [filterKey, filterValue] of Object.entries(this.containersFilters)) {
                filterValue = filterValue.trim()
                if (filterValue === '') {
                    continue
                }
                const filterKeyCapitalized = filterKey.charAt(0).toUpperCase() + filterKey.slice(1)
                queryParams.set('containers' + filterKeyCapitalized, filterValue)
            }

            htmx.ajax(
                'GET', '/overview/?' + queryParams.toString(),
                {
                    select: '#containers-table',
                    target: '#containers-table',
                    indicator: '#loading-overlay',
                    swap: 'outerHTML transition:true'
                },
            )
        },

        // Modal States
        isCreateContainerModalOpen: false,
        openCreateContainerModal() {
            this.resetPrimaryState()
            this.createContainerSelectedMarketplaceItemId = 0

            this.isCreateContainerModalOpen = true
        },
        closeCreateContainerModal() {
            this.isCreateContainerModalOpen = false
            this.$dispatch('clear:component-state')
        },
        isUpdateContainerModalOpen: false,
        openUpdateContainerModal(accountId, containerId, profileId) {
            this.resetPrimaryState()

            this.container.accountId = accountId
            this.container.id = containerId
            this.container.profileId = profileId
            this.isUpdateContainerModalOpen = true
        },
        closeUpdateContainerModal() {
            this.isUpdateContainerModalOpen = false
        },
        isShutdownContainerModalOpen: false,
        openShutdownContainerModal(accountId, containerId) {
            this.resetPrimaryState()

            this.container.accountId = accountId
            this.container.id = containerId
            this.isShutdownContainerModalOpen = true
        },
        updateContainerStatus(newStatus) {
            return htmx.ajax(
                'PUT', '/api/v1/container/',
                {
                    swap: 'none',
                    values: {
                        accountId: this.container.accountId,
                        containerId: this.container.id,
                        status: newStatus,
                    }
                },
            )
        },
        shutdownContainer() {
            this.updateContainerStatus(false).then(() => {
                this.$dispatch('update:container')
            })
            this.isShutdownContainerModalOpen = false
        },
        powerOnContainer(accountId, containerId) {
            this.resetPrimaryState()

            this.container.accountId = accountId
            this.container.id = containerId
            this.updateContainerStatus(true).then(() => {
                this.$dispatch('update:container')
            })
        },
        isRestartContainerModalOpen: false,
        openRestartContainerModal(accountId, containerId) {
            this.resetPrimaryState()

            this.container.accountId = accountId
            this.container.id = containerId
            this.isRestartContainerModalOpen = true
        },
        restartContainer() {
            this.updateContainerStatus(false).then(() => {
                this.updateContainerStatus(true).then(() => {
                    this.$dispatch('update:container')
                })
            })
            this.isRestartContainerModalOpen = false
        },
        isDeleteContainerModalOpen: false,
        openDeleteContainerModal(accountId, containerId, hostname) {
            this.resetPrimaryState()

            this.container.accountId = accountId
            this.container.id = containerId
            this.container.hostname = hostname
            this.isDeleteContainerModalOpen = true
        },
        closeDeleteContainerModal() {
            this.isDeleteContainerModalOpen = false
        },
        deleteContainer() {
            htmx.ajax(
                'DELETE',
                '/api/v1/container/' + this.container.accountId + '/' + this.container.id + '/',
                { swap: 'none' },
            ).then(() => {
                this.$dispatch('delete:container')
            })
            this.closeDeleteContainerModal()
        },
    }))

    Alpine.data('resourceUsage', () => ({
        // Primary State
        cpuMemoryChart: {
            chartObject: null,
            seriesNames: ['cpu', 'memory'],
            seriesObjects: {},
        },
        networkIoChart: {
            chartObject: null,
            seriesNames: ['rx', 'tx'],
            seriesObjects: {},
        },
        networkPacketsChart: {
            chartObject: null,
            seriesNames: ['vrx', 'vtx', 'drx', 'dtx'],
            seriesObjects: {},
        },
        networkErrorsChart: {
            chartObject: null,
            seriesNames: ['rx', 'tx'],
            seriesObjects: {},
        },
        storageIoChart: {
            chartObject: null,
            seriesNames: ['read', 'write'],
            seriesObjects: {},
        },
        storageIopsChart: {
            chartObject: null,
            seriesNames: ['read', 'write'],
            seriesObjects: {},
        },
        storageSpaceChart: {
            chartObject: null,
            seriesNames: ['used'],
            seriesObjects: {},
        },
        storageInodesChart: {
            chartObject: null,
            seriesNames: ['used', 'free'],
            seriesObjects: {},
        },

        // Auxiliary State
        refreshIntervalSecs: 20,
        sysInfo: {
            serverHostname: '',
            cpuModel: '',
            memoryRamGibibytes: 0,
            publicIpAddress: '',
            privateIpAddress: '',
            uptimeRelative: '',
        },
        networkChartTabSelected: 'io',
        storageChartTabSelected: 'io',
        chartOptions: {
            layout: {
                textColor: 'white',
                background: { type: 'solid', color: 'transparent' },
                fontFamily: "Lato, sans-serif",
                attributionLogo: false,
            },
            width: 350,
            height: 250,
            autoSize: false,
            timeScale: { timeVisible: true, barSpacing: 20 },
            grid: { vertLines: { color: '#ffffff1a' }, horzLines: { color: '#ffffff1a' } },
        },
        chartSeriesOptions: {
            0: { lineColor: '#6628d9', topColor: '#6628d9', bottomColor: 'rgb(46 16 101 / 30%)' },
            1: { lineColor: '#b31d3f', topColor: '#b31d3f', bottomColor: 'rgb(72 9 26 / 30%)' },
            2: { lineColor: '#a16207', topColor: '#a16207', bottomColor: 'rgb(66 32 6 / 30%)' },
            3: { lineColor: '#1d4ed8', topColor: '#1d4ed8', bottomColor: 'rgb(23 37 84 / 30%)' },
        },
        async updateResourceUsageCharts() {
            const o11yOverview = await fetch("/api/v1/o11y/overview/", {
                method: "GET",
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
            })
                .then((apiResponse) => {
                    if (!apiResponse.ok) {
                        throw new Error('BadHttpResponseCode: ' + apiResponse.status)
                    }

                    return apiResponse.json()
                })
                .then((jsonResponse) => {
                    return jsonResponse.body
                })
                .catch((error) => {
                    console.error(id + 'ReadO11yOverviewError: ', error)
                    return []
                })

            const nowUnixTime = Date.now() / 1000
            this.cpuMemoryChart.seriesObjects.cpu.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.cpuPercent },
            )
            this.cpuMemoryChart.seriesObjects.memory.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.memoryPercent },
            )
            let rxMebibytes = o11yOverview.resourceUsage.netInfoAggregated.recvBytes / 1024 / 1024
            if (rxMebibytes < 1) {
                rxMebibytes = 0
            }
            this.networkIoChart.seriesObjects.rx.update(
                { time: nowUnixTime, value: rxMebibytes },
            )
            let txMebibytes = o11yOverview.resourceUsage.netInfoAggregated.sentBytes / 1024 / 1024
            if (txMebibytes < 1) {
                txMebibytes = 0
            }
            this.networkIoChart.seriesObjects.tx.update(
                { time: nowUnixTime, value: txMebibytes },
            )
            this.networkPacketsChart.seriesObjects.vrx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.recvPackets },
            )
            this.networkPacketsChart.seriesObjects.vtx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.sentPackets },
            )
            this.networkPacketsChart.seriesObjects.drx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.recvDropPackets },
            )
            this.networkPacketsChart.seriesObjects.dtx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.sentDropPackets },
            )
            this.networkErrorsChart.seriesObjects.rx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.recvErrs },
            )
            this.networkErrorsChart.seriesObjects.tx.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.netInfoAggregated.sentErrs },
            )
            let readMebibytes = o11yOverview.resourceUsage.userDataStorageInfo.readBytes / 1024 / 1024
            if (readMebibytes < 1) {
                readMebibytes = 0
            }
            this.storageIoChart.seriesObjects.read.update(
                { time: nowUnixTime, value: readMebibytes },
            )
            let writeMebibytes = o11yOverview.resourceUsage.userDataStorageInfo.writeBytes / 1024 / 1024
            if (writeMebibytes < 1) {
                writeMebibytes = 0
            }
            this.storageIoChart.seriesObjects.write.update(
                { time: nowUnixTime, value: writeMebibytes },
            )
            this.storageIopsChart.seriesObjects.read.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.userDataStorageInfo.readOpsCount },
            )
            this.storageIopsChart.seriesObjects.write.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.userDataStorageInfo.writeOpsCount },
            )
            this.storageSpaceChart.seriesObjects.used.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.userDataStorageInfo.usedPercent },
            )
            this.storageInodesChart.seriesObjects.used.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.userDataStorageInfo.usedInodes },
            )
            this.storageInodesChart.seriesObjects.free.update(
                { time: nowUnixTime, value: o11yOverview.resourceUsage.userDataStorageInfo.freeInodes },
            )

            if (this.sysInfo.serverHostname != "") {
                return
            }
            this.sysInfo.serverHostname = o11yOverview.hostname
            this.sysInfo.cpuModel = o11yOverview.specs.cpuModelName + " // " +
                o11yOverview.specs.cpuCoresCount + " @ " + o11yOverview.specs.cpuFrequency + " Ghz"
            const memoryRamGibibytes = o11yOverview.specs.memoryTotalBytes / 1024 / 1024 / 1024
            this.sysInfo.memoryRamGibibytes = memoryRamGibibytes.toFixed(1) + " GiB"
            this.sysInfo.publicIpAddress = o11yOverview.publicIp
            this.sysInfo.privateIpAddress = o11yOverview.privateIp
            this.sysInfo.uptimeRelative = o11yOverview.uptimeRelative
        },

        init() {
            const chartNames = [
                'cpuMemoryChart', 'networkIoChart', 'networkPacketsChart',
                'networkErrorsChart', 'storageIoChart', 'storageIopsChart',
                'storageSpaceChart', 'storageInodesChart',
            ]
            for (const chartName of chartNames) {
                const domRef = this.$refs[chartName]
                if (domRef === undefined) {
                    continue
                }

                this[chartName].chartObject = LightweightCharts.createChart(domRef, this.chartOptions)
                for (const [seriesIndex, seriesName] of this[chartName].seriesNames.entries()) {
                    let seriesOptions = this.chartSeriesOptions[seriesIndex]
                    if (chartName === 'cpuMemoryChart' || chartName === 'storageSpaceChart') {
                        seriesOptions = Object.assign({}, seriesOptions)
                        seriesOptions.priceFormat = {
                            type: 'percent',
                            precision: 1,
                        }
                        seriesOptions.autoscaleInfoProvider = () => ({
                            priceRange: {
                                minValue: 0,
                                maxValue: 80,
                            },
                        })
                    }

                    this[chartName].seriesObjects[seriesName] = this[chartName].chartObject.addAreaSeries(
                        seriesOptions,
                    )
                }
            }

            setInterval(() => {
                this.updateResourceUsageCharts()
            }, parseInt(this.refreshIntervalSecs) * 1000)
        }
    }))
})