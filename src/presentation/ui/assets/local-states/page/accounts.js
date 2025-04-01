document.addEventListener("alpine:init", () => {
  Alpine.data("accounts", () => ({
    accountEntity: {},
    accountApiKey: "",
    resetPrimaryStates() {
      this.accountEntity = {
        id: "",
        username: "",
        password: "",
        quota: {
          cpuCores: 0,
          memoryGibibytes: 0,
          memoryMebibytes: 0,
          storageGibibytes: 0,
          storageMebibytes: 0,
          storageInodes: 0,
          storagePerformanceUnits: 1,
        },
      };
      this.accountApiKey = "";
    },
    updateAccountEntity(accountId) {
      this.accountEntity = JSON.parse(
        document.getElementById("accountEntity_" + accountId).textContent
      );
      delete this.accountEntity.password;
    },
    init() {
      this.resetPrimaryStates();
    },

    isMemoryGibibyteSelector: true,
    isStorageGibibyteSelection: true,

    isCreateAccountModalOpen: false,
    openCreateAccountModal() {
      this.resetPrimaryStates();
      this.isCreateAccountModalOpen = true;
    },
    closeCreateAccountModal() {
      this.isCreateAccountModalOpen = false;
    },
    isUpdateAccountQuotaModalOpen: false,
    openUpdateAccountQuotaModal(accountId) {
      this.resetPrimaryStates();
      this.updateAccountEntity(accountId);
      this.isUpdateAccountQuotaModalOpen = true;
    },
    closeUpdateAccountQuotaModal() {
      this.isUpdateAccountQuotaModalOpen = false;
    },
    isUpdatePasswordModalOpen: false,
    openUpdatePasswordModal(accountId) {
      this.resetPrimaryStates();
      this.updateAccountEntity(accountId);
      this.isUpdatePasswordModalOpen = true;
    },
    closeUpdatePasswordModal() {
      this.isUpdatePasswordModalOpen = false;
    },
    isUpdateApiKeyModalOpen: false,
    openUpdateApiKeyModal(accountId) {
      this.resetPrimaryStates();
      this.updateAccountEntity(accountId);
      this.isUpdateApiKeyModalOpen = true;
    },
    closeUpdateApiKeyModal() {
      this.isUpdateApiKeyModalOpen = false;
    },
    updateApiKey() {
      const shouldDisplayToast = false;
      Infinite.JsonAjax(
        "PUT",
        "/api/v1/account/",
        { id: this.accountEntity.id, shouldUpdateApiKey: true },
        shouldDisplayToast
      ).then((apiKey) => (this.accountApiKey = apiKey));
    },
    isDeleteAccountModalOpen: false,
    openDeleteAccountModal(accountId) {
      this.resetPrimaryStates();
      this.updateAccountEntity(accountId);
      this.isDeleteAccountModalOpen = true;
    },
    closeDeleteAccountModal() {
      this.isDeleteAccountModalOpen = false;
    },
    deleteAccount() {
      htmx
        .ajax("DELETE", "/api/v1/account/" + this.accountEntity.id + "/", {
          swap: "none",
        })
        .then(() => this.$dispatch("delete:account"));
      this.closeDeleteAccountModal();
    },
  }));
});
