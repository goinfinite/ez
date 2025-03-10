document.addEventListener("alpine:initializing", () => {
  Alpine.store("main", {
    textualViewSelector: Alpine.$persist(false).as("dash.textualViewSelector"),
    displayScheduledTasksPopover: Alpine.$persist(false).as(
      "dash.displayScheduledTasksPopover"
    ),

    toggleScheduledTasksPopover() {
      this.displayScheduledTasksPopover = !this.displayScheduledTasksPopover;
    },
    refreshScheduledTasksPopover() {
      window.dispatchEvent(new CustomEvent("refresh:footer"));
      setTimeout(() => {
        this.displayScheduledTasksPopover = true;
      }, 1000);
    },

    init() {
      Alpine.watch(
        () => this.textualViewSelector,
        () => {
          htmx.process(document.body);
        }
      );
    },
  });
});
