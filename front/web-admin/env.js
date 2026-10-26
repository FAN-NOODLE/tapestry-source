const hash = process.argv[2];
const apiDomain = process.argv[3];
const mode = process.argv[4];
const fs = require("fs");
const string =
  "NODE_ENV = production \n" +
  "VUE_APP_HASH = " +
  hash +
  "\n" +
  "VUE_APP_API_ROOT = " +
  apiDomain +
  "\n" +
  "VUE_APP_MODE =" +
  mode +
  "\n" +
  "VUE_APP_DEVOPS = old";
fs.writeFile(".env.production", string, function(err) {
  if (err) {
    return console.log(err);
  }
  console.log("The file was saved!");
});
