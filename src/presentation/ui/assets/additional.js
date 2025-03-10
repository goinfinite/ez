"use strict";

// UnoCSS customizations
window.__unocss = {
  theme: {
    colors: {
      infinite: {
        50: "#dea893",
        100: "#d89a81",
        200: "#d38b6f",
        300: "#cd7d5d",
        400: "#ca7654",
        500: "#c97350",
        600: "#c46f4d",
        700: "#ba6949",
        800: "#a55e41",
        900: "#905239",
        950: "#7c4631",
      },
      ez: {
        50: "#667b73",
        100: "#4c655c",
        200: "#334f45",
        300: "#19392d",
        400: "#0c2e22",
        500: "#06261b",
        600: "#06271b",
        700: "#06251a",
        800: "#052117",
        900: "#041d14",
        950: "#041911",
      },
      os: {
        5: "#66737a",
        100: "#4d5b64",
        200: "#34444e",
        300: "#1a2d38",
        400: "#0d212d",
        500: "#071b27",
        600: "#071a26",
        700: "#061924",
        800: "#061620",
        900: "#05131c",
        950: "#041118",
      },
    },
  },
};

document.addEventListener("alpine:init", () => {
  async function jsonAjax(method, url, payload, shouldDisplayToast = true) {
    const loadingOverlayElement = document.getElementById("loading-overlay");
    loadingOverlayElement.classList.add("htmx-request");

    try {
      const requestSettings = {
        method: method,
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };
      if (Object.keys(payload).length > 0) {
        requestSettings.body = JSON.stringify(payload);
      }
      const response = await fetch(url, requestSettings);
      const parsedResponse = await response.json();

      loadingOverlayElement.classList.remove("htmx-request");

      if (!response.ok) {
        throw new Error(parsedResponse.body);
      }

      if (shouldDisplayToast) {
        Alpine.store("toast").displayToast(parsedResponse.body, "success");
      }
      return parsedResponse.body;
    } catch (error) {
      loadingOverlayElement.classList.remove("htmx-request");

      if (shouldDisplayToast) {
        Alpine.store("toast").displayToast(error.message, "danger");
        return;
      }
      throw error;
    }
  }

  function createRandomPassword() {
    const passwordLength = 16;
    const chars =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+";

    let passwordContent = "";
    let passwordIterationCount = 0;
    while (passwordIterationCount < passwordLength) {
      const randomIndex = Math.floor(Math.random() * chars.length);
      const indexAfterRandomIndex = randomIndex + 1;
      passwordContent += chars.substring(randomIndex, indexAfterRandomIndex);

      passwordIterationCount++;
    }

    return passwordContent;
  }

  window.Infinite = {
    JsonAjax: jsonAjax,
    CreateRandomPassword: createRandomPassword,
  };
});
