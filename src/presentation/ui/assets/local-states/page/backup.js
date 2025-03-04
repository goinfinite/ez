document.addEventListener("alpine:init", () => {
  Alpine.data("backupIndex", () => ({
    // Primary States
    backupFeatureTabSelected: "tasks",
  }));

  Alpine.data("backupTasks", () => ({
    // Primary State
    taskEntity: {},
    createTaskArchive: {},
    restoreTask: {},
    resetPrimaryState() {
      this.taskEntity = {};
      this.createTaskArchive = {};
      this.restoreTask = {};
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
      this.resetPrimaryState();
      this.updateTaskEntity(taskId);
      this.createTaskArchive = { taskId: taskId };
      this.isCreateTaskArchiveModalOpen = true;
    },
    closeCreateTaskArchiveModal() {
      this.isCreateTaskArchiveModalOpen = false;
    },

    isRestoreTaskModalOpen: false,
    openRestoreTaskModal(taskId) {
      this.resetPrimaryState();
      this.updateTaskEntity(taskId);
      this.restoreTask = { taskId: taskId };
      this.isRestoreTaskModalOpen = true;
    },
    closeRestoreTaskModal() {
      this.isRestoreTaskModalOpen = false;
    },
  }));
});
