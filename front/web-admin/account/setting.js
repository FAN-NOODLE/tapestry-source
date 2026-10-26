/* eslint-disable @typescript-eslint/no-var-requires */
/**
 * @description 本地运行测试环境 npm run test
 * @description 本地运行生产环境 npm run prod (代理解决跨域)
 * */
require("colors");
const inquirer = require("inquirer");
const devConfig = require("./dev.conf");
const testConfig = require("./test.conf");
const prodConfig = require("./prod.conf");
const fs = require("fs");
const env = process.argv[2];
const envMap = {
  dev: { config: devConfig, name: "开发" },
  test: { config: testConfig, name: "测试" },
  prod: { config: prodConfig, name: "生产" }
};
const config = envMap[env].config;
const promptList = [
  {
    type: "list",
    message: `请选择要运行的${envMap[env].name}环境:`,
    name: "type",
    choices: Object.keys(config)
  }
];

const setAcount = ({ account, pswd }) => {
  const file = fs.readFileSync("src/utils/LoginManager.ts", "utf8");
  let resultFile = file.replace(
    /MOCK_ACCOUNT = ([\s\S]*?);/g,
    `MOCK_ACCOUNT = "${account}";`
  );
  resultFile = resultFile.replace(
    /PASS_WORD = ([\s\S]*?);/g,
    `PASS_WORD = "${pswd}";`
  );
  fs.writeFileSync("src/utils/LoginManager.ts", resultFile, "utf8");
  console.log(`切换${envMap[env].name}环境成功,开始运行`.green);
};
(async function() {
  const data = await inquirer.prompt(promptList);
  if (config[data.type]) {
    setAcount(config[data.type]);
  } else {
    console.log("没找到配置的企业账号密码，请检查文件".red);
  }
})();
