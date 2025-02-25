document.addEventListener("alpine:init", () => {
  Alpine.data("backupIndex", () => ({
    // Primary States
    backupFeatureTabSelected: "tasks",
  }));
});
