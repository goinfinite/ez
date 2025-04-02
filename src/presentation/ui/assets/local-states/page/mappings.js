document.addEventListener("alpine:init", () => {
    Alpine.data("mappings", () => ({
        isCreateMappingModalOpen: false,
        openCreateMappingModal() {
            this.isCreateMappingModalOpen = true
        },
        closeCreateMappingModal() {
            this.isCreateMappingModalOpen = false
        },
    }))
})