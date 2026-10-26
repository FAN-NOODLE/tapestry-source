let VUE_APP_API_ROOT = "";

if (window.location.hostname === "tapestry.vchangyi.com") {
  VUE_APP_API_ROOT = "https://tapestry.vchangyi.com/";
} else {
  VUE_APP_API_ROOT = "https://tapestry-dev.vchangyi.com/";
}

export { VUE_APP_API_ROOT };
