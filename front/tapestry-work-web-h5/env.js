let hash = process.argv[2];
let apiDomain = process.argv[3];
const fs = require("fs");
let string =
  "NODE_ENV = production \n" +
  "VUE_APP_HASH = " +
  hash +
  "\n" +
  "VUE_APP_API_ROOT = " +
  apiDomain;
fs.writeFile(".env.production", string, function (err) {
  if (err) {
    return console.log(err);
  }
  console.log("The file was saved!");
});
