document.addEventListener('alpine:init', () => {
    Alpine.data('imageArchive', () => ({
        // Primary State
        imageArchive: {},
        resetPrimaryState() {
            this.imageArchive = {
                'id': '',
                'accountId': '',
            }
        },
        init() {
            this.resetPrimaryState()
        },

        // Modal States
        isDeleteArchiveModalOpen: false,
        openDeleteArchiveModal(accountId, imageId) {
            this.resetPrimaryState()

            this.imageArchive.accountId = accountId
            this.imageArchive.id = imageId
            this.isDeleteArchiveModalOpen = true
        },
        closeDeleteArchiveModal() {
            this.isDeleteArchiveModalOpen = false
        },
        deleteArchive() {
            htmx.ajax(
                'DELETE',
                '/api/v1/container/image/archive/' + this.imageArchive.accountId + '/' + this.imageArchive.id + '/',
                { swap: 'none' },
            )
            this.$dispatch('delete:image-archive')
            this.closeDeleteArchiveModal()
        },
    }))
    Alpine.data('containerImage', () => ({
        // Primary State
        containerImage: {},
        resetPrimaryState() {
            this.containerImage = {
                'id': '',
                'accountId': '',
                'imageAddress': '',
                'imageHash': '',
                'isa': '',
                'sizeBytes': '',
                'portBindings': '',
                'envs': '',
                'createdAt': '',
            }
        },
        init() {
            this.resetPrimaryState()
        },

        // Auxiliary States
        archiveImage(accountId, imageId) {
            htmx.ajax(
                'POST', '/api/v1/container/image/archive/',
                { swap: 'none', values: { accountId: accountId, imageId: imageId } },
            )
            this.$store.main.refreshScheduledTasksPopover()
        },

        // Modal States
        isImportImageModalOpen: false,
        openImportImageModal() {
            this.resetPrimaryState()

            this.isImportImageModalOpen = true
        },
        closeImportImageModal() {
            this.isImportImageModalOpen = false
        },
        isTakeSnapshotModalOpen: false,
        openTakeSnapshotModal() {
            this.resetPrimaryState()

            this.isTakeSnapshotModalOpen = true
        },
        closeTakeSnapshotModal() {
            this.isTakeSnapshotModalOpen = false
        },
        isDeleteImageModalOpen: false,
        openDeleteImageModal(accountId, imageId) {
            this.resetPrimaryState()

            this.containerImage.accountId = accountId
            this.containerImage.id = imageId
            this.isDeleteImageModalOpen = true
        },
        closeDeleteImageModal() {
            this.isDeleteImageModalOpen = false
        },
        deleteImage() {
            htmx.ajax(
                'DELETE',
                '/api/v1/container/image/' + this.containerImage.accountId + '/' + this.containerImage.id + '/',
                { swap: 'none' },
            )
            this.$dispatch('delete:container-image')
            this.closeDeleteImageModal()
        },
    }))
    Alpine.data('containerSnapshot', () => ({
        // Primary State
        accountId: '',
        containerId: '',
        shouldCreateArchive: false,
        shouldDiscardImage: false,
        archiveCompressionFormat: 'br',
        resetPrimaryState() {
            this.accountId = ''
            this.containerId = ''
            this.shouldCreateArchive = false
            this.shouldDiscardImage = false
            this.archiveCompressionFormat = 'br'
        },
        init() {
            this.resetPrimaryState()
        },
    }))
})