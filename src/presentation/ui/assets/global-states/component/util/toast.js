document.addEventListener("alpine:init", () => {
  Alpine.store("toast", {
    toastVisible: false,
    toastMessage: "",
    toastType: "danger",

    displayToast(message, type) {
      this.toastVisible = true;
      this.toastMessage = message;
      this.toastType = type;
      setTimeout(() => {
        this.clearToast();
      }, 3000);
    },

    clearToast() {
      this.toastVisible = false;
      this.toastMessage = "";
    },
  });
});

document.addEventListener("htmx:afterRequest", (event) => {
  const contentType = event.detail.xhr.getResponseHeader("Content-Type");
  if (contentType !== "application/json") {
    return;
  }

  const responseData = event.detail.xhr.responseText;
  if (responseData === "") {
    return;
  }

  let toastType = "success";
  const isResponseError = event.detail.xhr.status >= 400;
  if (isResponseError) {
    toastType = "danger";
  }

  const parsedResponse = JSON.parse(responseData);
  if (parsedResponse.body === undefined || parsedResponse.body === "") {
    return;
  }
  if (typeof parsedResponse.body !== "string") {
    return;
  }

  const toastMessage = parsedResponse.body;

  Alpine.store("toast").displayToast(toastMessage, toastType);
});
