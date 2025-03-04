function InitReadRequestSettingsLocalState(
  componentId,
  uiEndpoint,
  recordsDisplayElementId
) {
  document.addEventListener("alpine:init", () => {
    Alpine.data(componentId + "ReadRequestSettings", () => ({
      pagination: {},
      filters: {},

      resetState() {
        const readRequestDto = JSON.parse(
          document.getElementById(componentId + "ReadRequestDto").textContent
        );
        this.filters = Object.assign({}, readRequestDto);
        delete this.filters.pagination;
        for (let [filterKey, filterValue] of Object.entries(this.filters)) {
          if (filterValue === null) {
            this.filters[filterKey] = "";
          }
        }

        this.pagination = JSON.parse(
          document.getElementById(componentId + "PaginationDto").textContent
        );
      },

      readQueryParams() {
        queryParams = new URLSearchParams();
        queryParams.set(componentId + "PageNumber", this.pagination.pageNumber);
        queryParams.set(
          componentId + "ItemsPerPage",
          this.pagination.itemsPerPage
        );

        for (let [filterKey, filterValue] of Object.entries(this.filters)) {
          if (filterValue === null) {
            continue;
          }

          if (
            typeof filterValue === "string" ||
            filterValue instanceof String
          ) {
            filterValue = filterValue.trim();
          }
          if (filterValue === "") {
            continue;
          }
          const filterKeyPascalCase =
            filterKey.charAt(0).toUpperCase() + filterKey.slice(1);
          queryParams.set(componentId + filterKeyPascalCase, filterValue);
        }

        return queryParams.toString();
      },

      reloadRecordsDisplay() {
        htmx.ajax("GET", uiEndpoint + "?" + this.readQueryParams(), {
          select: recordsDisplayElementId,
          target: recordsDisplayElementId,
          indicator: "#loading-overlay",
          swap: "outerHTML transition:true",
        });
      },

      init() {
        this.resetState();
        this.$watch("pagination", () => {
          this.reloadRecordsDisplay();
        });
      },
    }));
  });
}
