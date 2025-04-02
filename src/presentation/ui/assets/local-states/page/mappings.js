document.addEventListener("alpine:init", () => {
    Alpine.data("mappings", () => ({
        mappingEntity: {},
        targeEntity: {},
        resetPrimaryStates() {
            this.mappingEntity = { id: "", hostname: "" }
            this.targetEntity = { id: "" }
        },
        init() {
            this.resetPrimaryStates()
        },

        isCreateMappingModalOpen: false,
        openCreateMappingModal() {
            this.resetPrimaryStates()
            this.isCreateMappingModalOpen = true
        },
        closeCreateMappingModal() {
            this.isCreateMappingModalOpen = false
        },
        isDeleteMappingModalOpen: false,
        openDeleteMappingModal(mappingId, mappingHostname) {
            this.resetPrimaryStates()
            this.mappingEntity.id = mappingId
            this.mappingEntity.hostname = mappingHostname
            this.isDeleteMappingModalOpen = true
        },
        closeDeleteMappingModal() {
            this.isDeleteMappingModalOpen = false
        },
        deleteMapping() {
            htmx
                .ajax("DELETE", "/api/v1/mapping/" + this.mappingEntity.id + "/", {
                    swap: "none",
                })
                .then(() => this.$dispatch("refresh:mappings-table"))
            this.closeDeleteMappingModal()
        },
        isCreateTargetModalOpen: false,
        openCreateTargetModal(mappingId = "", mappingHostname = "") {
            this.resetPrimaryStates()
            this.mappingEntity.id = mappingId
            this.mappingEntity.hostname = mappingHostname
            this.isCreateTargetModalOpen = true
        },
        closeCreateTargetModal() {
            this.isCreateTargetModalOpen = false
        },
        isDeleteTargetModalOpen: false,
        openDeleteTargetModal(mappingId, targetId) {
            this.resetPrimaryStates()
            this.mappingEntity.id = mappingId
            this.targetEntity.id = targetId
            this.isDeleteTargetModalOpen = true
        },
        closeDeleteTargetModal() {
            this.isDeleteTargetModalOpen = false
        },
        deleteTarget() {
            htmx
                .ajax("DELETE", "/api/v1/mapping/" + this.mappingEntity.id + "/target/" + this.targetEntity.id + "/", {
                    swap: "none",
                })
                .then(() => this.$dispatch("refresh:mappings-table"))
            this.closeDeleteTargetModal()
        },
    }))
})