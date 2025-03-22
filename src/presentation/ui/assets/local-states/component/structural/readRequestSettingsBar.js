function InitReadRequestSettingsLocalState(
  componentId,
  uiEndpoint,
  recordsDisplayElementId
) {
  document.addEventListener("alpine:init", () => {
    Alpine.data(componentId + "ReadRequestSettings", () => ({
      pagination: {},
      filters: {},
      entityFieldKeys: [],

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
        const entityFields = JSON.parse(
          document.getElementById(componentId + "Struct").textContent
        );
        this.entityFieldKeys = Object.keys(entityFields);
      },

      readQueryParams() {
        queryParams = new URLSearchParams();

        let paginationWithFilters = Object.assign(
          {},
          this.pagination,
          this.filters
        );
        delete paginationWithFilters.pagesTotal;
        delete paginationWithFilters.itemsTotal;

        for (let [fieldKey, fieldValue] of Object.entries(
          paginationWithFilters
        )) {
          if (fieldValue === null) {
            continue;
          }

          if (typeof fieldValue === "string" || fieldValue instanceof String) {
            fieldValue = fieldValue.trim();
          }
          if (fieldValue === "") {
            continue;
          }
          const fieldKeyPascalCase =
            fieldKey.charAt(0).toUpperCase() + fieldKey.slice(1);
          queryParams.set(componentId + fieldKeyPascalCase, fieldValue);
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
