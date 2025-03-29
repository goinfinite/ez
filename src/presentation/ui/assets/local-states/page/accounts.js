document.addEventListener('alpine:init', () => {
    Alpine.data('accounts', () => ({
        accountEntity: {},
        resetPrimaryStates() {
            this.accountEntity = {}
        },
        updateAccountEntity(accountId) {
            this.accountEntity = JSON.parse(
                document.getElementById("accountEntity_" + accountId).textContent
            )
        },
        init() {
            this.resetPrimaryStates()
        },

        isUpdatePasswordModalOpen: false,
        openUpdatePasswordModal(accountId) {
            this.resetPrimaryStates()
            this.updateAccountEntity(accountId)
            this.isUpdatePasswordModalOpen = true
        },
        closeUpdatePasswordModal() {
            this.isUpdatePasswordModalOpen = false
        },
        isDeleteAccountModalOpen: false,
        openDeleteAccountModal(accountId) {
            this.resetPrimaryStates()
            this.updateAccountEntity(accountId)
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
