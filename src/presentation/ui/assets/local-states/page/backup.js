document.addEventListener("alpine:init", () => {
  Alpine.data("backupIndex", () => ({
    // Primary States
    backupFeatureTabSelected: "tasks",
  }));

  Alpine.data("backupTasks", () => ({
    // Primary State
    taskEntity: {},
    createTaskArchive: {},
    resetPrimaryState() {
      this.taskEntity = {};
      this.createTaskArchive = {};
    },
    updateTaskEntity(taskId) {
      this.taskEntity = JSON.parse(
        document.getElementById("backupTaskEntity_" + taskId).textContent
      );
    },
    init() {
      this.resetPrimaryState();
    },

    // Auxiliary State
    isCreateTaskArchiveModalOpen: false,

    openCreateTaskArchiveModal(taskId) {
      this.updateTaskEntity(taskId);
      this.createTaskArchive = {
        taskId: taskId,
        containerAccountIds: [],
      };
      this.isCreateTaskArchiveModalOpen = true;
    },
    closeCreateTaskArchiveModal() {
      this.resetPrimaryState();
      this.isCreateTaskArchiveModalOpen = false;
    },
  }));
});
