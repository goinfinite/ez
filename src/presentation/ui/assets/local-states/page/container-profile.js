document.addEventListener('alpine:init', () => {
	Alpine.data('containerProfile', () => ({
		// Primary State
		containerProfile: {},
		resetPrimaryState() {
			this.containerProfile = {
				'id': '',
				'accountId': '',
				'name': '',
				'baseSpecs': {
					'millicores': 0,
					'cpuCores': 0,
					'memoryBytes': 0,
					'memoryMebibytes': 0,
					'memoryGibibytes': 0,
					'storagePerformanceUnits': 0
				},
				'maxSpecs': {
					'millicores': null,
					'cpuCores': null,
					'memoryBytes': null,
					'memoryMebibytes': null,
					'memoryGibibytes': null,
					'storagePerformanceUnits': null
				},
				'scalingPolicy': null,
				'scalingThreshold': null,
				'scalingMaxDurationSecs': null,
				'scalingMaxDurationMinutes': null,
				'scalingMaxDurationHours': null,
				'scalingIntervalSecs': null,
				'scalingIntervalMinutes': null,
				'scalingIntervalHours': null,
				'hostMinCapacityPercent': null
			}
		},
		init() {
			this.resetPrimaryState()
		},

		// Auxiliary States
		gibibyteSelector: true,
		preferedByteSuffix: 'GiB',
		scalingIntervalHoursSelector: true,
		scalingMaxDurationHoursSelector: true,
		preferedScalingIntervalSuffix: 'hour(s)',
		preferedScalingMaxDurationSuffix: 'hour(s)',
		get isScalingPolicyConn() {
			return this.containerProfile.scalingPolicy == 'connection'
		},
		resetAuxiliaryStates() {
			this.gibibyteSelector = true
			this.preferedByteSuffix = 'GiB'
			this.scalingIntervalHoursSelector = true
			this.scalingMaxDurationHoursSelector = true
			this.preferedScalingIntervalSuffix = 'hour(s)'
			this.preferedScalingMaxDurationSuffix = 'hour(s)'
		},

		// Modal States
		isUpdateModalOpen: false,
		openUpdateModal(containerProfileData) {
			this.resetPrimaryState()
			this.resetAuxiliaryStates()

			this.containerProfile = containerProfileData
			if (this.containerProfile.baseSpecs.memoryGibibytes < 1) {
				this.preferedByteSuffix = 'MiB'
				this.gibibyteSelector = false
			}

			if (this.containerProfile.scalingIntervalHours < 1) {
				this.scalingIntervalHoursSelector = false
				this.preferedScalingIntervalSuffix = 'min(s)'
			}

			if (this.containerProfile.scalingMaxDurationHours < 1) {
				this.scalingMaxDurationHoursSelector = false
				this.preferedScalingMaxDurationSuffix = 'min(s)'
			}

			this.isUpdateModalOpen = true
		},
		closeUpdateModal() {
			this.isUpdateModalOpen = false
		},
		isCreateModalOpen: false,
		openCreateModal() {
			this.resetPrimaryState()
			this.resetAuxiliaryStates()

			this.isCreateModalOpen = true
		},
		closeCreateModal() {
			this.isCreateModalOpen = false
		},
		isDeleteModalOpen: false,
		openDeleteModal(accountId, profileId, profileName) {
			this.resetPrimaryState()
			this.resetAuxiliaryStates()

			this.containerProfile.id = profileId
			this.containerProfile.accountId = accountId
			this.containerProfile.name = profileName
			this.isDeleteModalOpen = true
		},
		closeDeleteModal() {
			this.isDeleteModalOpen = false
		},
		deleteElement() {
			htmx.ajax(
				'DELETE',
				'/api/v1/container/profile/' + this.containerProfile.accountId + '/' + this.containerProfile.id + '/',
				{ swap: 'none' },
			)
			this.closeDeleteModal()
		},
	}))
})