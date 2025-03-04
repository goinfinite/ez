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

    isDeleteTaskModalOpen: false,
    openDeleteTaskModal(taskId) {
      this.resetPrimaryState();
      this.updateTaskEntity(taskId);
      this.isDeleteTaskModalOpen = true;
    },
    closeDeleteTaskModal() {
      this.isDeleteTaskModalOpen = false;
    },
    deleteTask() {
      htmx
        .ajax("DELETE", "/api/v1/backup/tasks/" + this.taskEntity.id + "/", {
          swap: "none",
        })
        .then(() => {
          this.$dispatch("delete:backup-task");
        });
      this.closeDeleteTaskModal();
    },
  }));

  Alpine.data("backupTaskArchives", () => ({
    taskArchiveEntity: {},
    restoreTaskArchive: {},
    resetPrimaryState() {
      this.taskArchiveEntity = {};
      this.restoreTaskArchive = {};
    },
    updateTaskArchiveEntity(taskArchiveId) {
      this.taskArchiveEntity = JSON.parse(
        document.getElementById("backupTaskArchiveEntity_" + taskArchiveId)
          .textContent
      );
    },
    init() {
      this.resetPrimaryState();
    },

    isRestoreTaskArchiveModalOpen: false,
    openRestoreTaskArchiveModal(archiveId) {
      this.resetPrimaryState();
      this.updateTaskArchiveEntity(archiveId);
      this.restoreTaskArchive = { archiveId: archiveId };
      this.isRestoreTaskArchiveModalOpen = true;
    },
    closeRestoreTaskArchiveModal() {
      this.isRestoreTaskArchiveModalOpen = false;
    },

    isDeleteTaskArchiveModalOpen: false,
    openDeleteTaskArchiveModal(taskArchiveId) {
      this.resetPrimaryState();
      this.updateTaskArchiveEntity(taskArchiveId);
      this.isDeleteTaskArchiveModalOpen = true;
    },
    closeDeleteTaskArchiveModal() {
      this.isDeleteTaskArchiveModalOpen = false;
    },
    deleteTaskArchive() {
      htmx
        .ajax(
          "DELETE",
          "/api/v1/backup/task/archive/" + this.taskArchiveEntity.id + "/",
          { swap: "none" }
        )
        .then(() => {
          this.$dispatch("delete:backup-task-archive");
        });
      this.closeDeleteTaskArchiveModal();
    },
  }));
});
