document.addEventListener('alpine:init', () => {
    Alpine.data('accounts', () => ({
        // Primary States
        accountEntity: {},
        resetPrimaryStates() {
            this.accountEntity = {}
        },
        init() {
            this.resetPrimaryStates()
        },

        // Modal States
        isDeleteAccountModalOpen: false,
        openDeleteAccountModal(accountId) {
            this.resetPrimaryStates()

            this.accountEntity = JSON.parse(
                document.getElementById("accountEntity_" + accountId).textContent
            )

            this.isDeleteAccountModalOpen = true
        },
        closeDeleteAccountModal() {
            this.isDeleteAccountModalOpen = false
        },
        deleteAccount() {
            htmx.ajax(
                "DELETE",
                "/api/v1/account/" + this.accountEntity.id + "/",
                { swap: "none" },
            ).then(() => this.$dispatch("delete:account"))
            this.closeDeleteAccountModal()
        },
    }))
})
