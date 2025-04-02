document.addEventListener("alpine:init", () => {
    Alpine.data("mappings", () => ({
        mappingEntity: {},
        resetPrimaryStates() {
            this.mappingEntity = { id: "", hostname: "" }
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
    }))
})