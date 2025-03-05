document.addEventListener("alpine:init", () => {
  backupApiBaseEndpoint = "/api/v1/backup";

  Alpine.data("backupIndex", () => ({
    // Primary States
    backupFeatureTabSelected: "jobs",
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
        .ajax(
          "DELETE",
          backupApiBaseEndpoint + "/task/" + this.taskEntity.taskId + "/",
          {
            swap: "none",
          }
        )
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
          backupApiBaseEndpoint +
            "/task/archive/" +
            this.taskArchiveEntity.archiveId +
            "/",
          { swap: "none" }
        )
        .then(() => {
          this.$dispatch("delete:backup-task-archive");
        });
      this.closeDeleteTaskArchiveModal();
    },
  }));

  Alpine.data("backupJobs", () => ({
    // Primary State
    jobEntity: {},
    createJob: {},
    resetPrimaryState() {
      this.jobEntity = {};
      this.createJob = {};
    },
    updateJobEntity(jobId) {
      this.jobEntity = JSON.parse(
        document.getElementById("backupJobEntity_" + jobId).textContent
      );
    },
    init() {
      this.resetPrimaryState();
    },

    // Auxiliary State
    jobApiEndpoint: backupApiBaseEndpoint + "/job",
    isRunJobModalOpen: false,
    openRunJobModal(jobId) {
      this.resetPrimaryState();
      this.updateJobEntity(jobId);
      this.isRunJobModalOpen = true;
    },
    closeRunJobModal() {
      this.isRunJobModalOpen = false;
    },
    runJob() {
      htmx.ajax(
        "POST",
        this.jobApiEndpoint +
          "/" +
          this.jobEntity.accountId +
          "/" +
          this.jobEntity.jobId +
          "/run/",
        {
          swap: "none",
        }
      );
      this.closeRunJobModal();
    },

    isCreateJobModalOpen: false,
    openCreateJobModal() {
      this.resetPrimaryState();
      this.isCreateJobModalOpen = true;
    },
    closeCreateJobModal() {
      this.isCreateJobModalOpen = false;
    },

    isUpdateJobModalOpen: false,
    openUpdateJobModal(jobId) {
      this.resetPrimaryState();
      this.updateJobEntity(jobId);
      this.isUpdateJobModalOpen = true;
    },
    closeUpdateJobModal() {
      this.isUpdateJobModalOpen = false;
    },

    isDeleteJobModalOpen: false,
    openDeleteJobModal(jobId) {
      this.resetPrimaryState();
      this.updateJobEntity(jobId);
      this.isDeleteJobModalOpen = true;
    },
    closeDeleteJobModal() {
      this.isDeleteJobModalOpen = false;
    },
    deleteJob() {
      htmx
        .ajax(
          "DELETE",
          this.jobApiEndpoint +
            "/" +
            this.jobEntity.accountId +
            "/" +
            this.jobEntity.jobId +
            "/",
          {
            swap: "none",
          }
        )
        .then(() => {
          this.$dispatch("delete:backup-job");
        });
      this.closeDeleteJobModal();
    },
  }));
});
