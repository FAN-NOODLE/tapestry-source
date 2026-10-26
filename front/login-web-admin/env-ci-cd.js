const mode = process.argv[2];
const fs = require("fs");

const apiConfig = {
  dev: "https://crs-api-dev.vchangyi.com/",
  test: "https://crs-api-test.vchangyi.com/",
  prod: "https://crs-api.vchangyi.com/"
};

const string =
  "NODE_ENV = production \n" +
  "VUE_APP_API_ROOT = " +
  apiConfig[mode] +
  "\n" +
  "VUE_APP_MODE =" +
  mode +
  "\n" +
  "VUE_APP_DEVOPS = new";
fs.writeFile(".env.production", string, function(err) {
  if (err) {
    return console.log(err);
  }
  console.log("The file was saved!");
});
