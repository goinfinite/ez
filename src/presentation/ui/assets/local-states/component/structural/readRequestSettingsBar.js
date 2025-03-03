function InitReadRequestSettingsLocalState(
  componentId,
  uiEndpoint,
  recordsDisplayElementId
) {
  document.addEventListener("alpine:init", () => {
    Alpine.data(componentId + "ReadRequestSettings", () => ({
      pagination: {},
      filters: {},

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

          filterValue = filterValue.trim();
          if (filterValue === "") {
            continue;
          }
          const filterKeyPascalCase =
            filterKey.charAt(0).toUpperCase() + filterKey.slice(1);
          queryParams.set(componentId + filterKeyPascalCase, filterValue);
        }

        return queryParams.toString();
      },

      reloadRecords() {
        htmx.ajax("GET", uiEndpoint + "?" + this.readQueryParams(), {
          select: recordsDisplayElementId,
          target: recordsDisplayElementId,
          indicator: "#loading-overlay",
          swap: "outerHTML transition:true",
        });
      },

      init() {
        const readRequestDto = JSON.parse(
          document.getElementById(componentId + "ReadRequestDto").textContent
        );
        this.pagination = Object.assign({}, readRequestDto.pagination);
        this.pagination.itemsPerPage = 5;
        this.filters = Object.assign({}, readRequestDto);
        for (let [filterKey, filterValue] of Object.entries(this.filters)) {
          if (filterValue === null) {
            this.filters[filterKey] = "";
          }
        }

        delete this.filters.pagination;
        document.addEventListener(
          "update:" + componentId + "-pagination",
          () => {
            this.reloadRecords();
          }
        );
      },
    }));
  });
}
