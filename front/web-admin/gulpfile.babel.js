"use strict";

const gulp = require("gulp");
const path = require("path");
const rename = require("gulp-rename");
const template = require("gulp-template");
const yargs = require("yargs");

const root = "src";

//创建页面 解析路径辅方法
const resolveToPages = (glob = "") => {
  return path.join(root, "views/myapp/", glob); // views/myapp/{glob}
};

//创建demo页面 解析路径辅方法
const resolveToDemoPages = (glob = "") => {
  return path.join(root, "views/demo/", glob); // views/myapp/{glob}
};

// 模版页面路径
const paths = {
  blankTemplates: path.join(__dirname, "generator", "temp/**/*.**")
};

gulp.task("page", () => {
  const cap = val => {
    return val.charAt(0).toUpperCase() + val.slice(1);
  };
  const calModulePath = path => {
    path = path.replace(/\./g, "/");
    path = path.replace(/\/+/g, "/");
    path = path.charAt(0) === "/" ? path.substring(1) : path;
    path =
      path.charAt(path.length - 1) === "/"
        ? path.substring(0, path.length - 1)
        : path;
    return path === "" ? path : "/" + path;
  };
  const moduleName = path => {
    const modulesConfig = [
      "app",
      "customer",
      "enterprise",
      "goods",
      "home",
      "material",
      "order",
      "setting",
      "shop",
      "store"
    ];
    let moduleName = "";
    modulesConfig.forEach(item => {
      if (path.indexOf("/" + item + "/") !== -1) {
        moduleName = item;
      }
    });
    return moduleName;
  };
  const componentInfo = yargs.argv.fileName.split(":"); // 冒号用来确认文件名和路径 冒号前为要创建的文件名，冒号后为此文件的路径 此规则默认为畅移所有项目
  const name = componentInfo[0]; // 默认为文件名
  const parentPath = (componentInfo.length > 1 && componentInfo[1]) || ""; // 判断是否需要生成路径
  const fileName = name.replace(/([A-Z])/g, "-$1").toLowerCase(); // 文件名称转化为小写
  const destPath = path.join(resolveToPages(), parentPath, fileName); // 文件生成的整体路径

  const nameOk = name.replace(/-(\w)/g, function(all, letter) {
    return letter.toUpperCase();
  }); // 文件名称转化为大驼峰
  // console.log(
  //   "文件名称小驼峰",
  //   nameOk,
  //   "文件名默认",
  //   fileName,
  //   "文件名大驼峰",
  //   cap(nameOk),
  //   "指定创建的文件路径",
  //   parentPath,
  //   "指定创建的文件路径加'/'",
  //   calModulePath(parentPath),
  //   "最终文件路径",
  //   destPath,
  //   "模块名称",
  //   moduleName(destPath)
  // );
  return gulp
    .src(paths.blankTemplates)
    .pipe(
      template({
        name: nameOk, //文件名称小驼峰
        fileName: fileName, //文件名默认
        upCaseName: cap(nameOk), // 文件名大驼峰
        parentPath: parentPath,
        modulePath: calModulePath(parentPath),
        moduleName: moduleName(destPath)
      })
    )
    .pipe(
      rename(path => {
        path.basename = path.basename.replace("temp", fileName);
      })
    )
    .pipe(gulp.dest(destPath));
});

gulp.task("demo", () => {
  const cap = val => {
    return val.charAt(0).toUpperCase() + val.slice(1);
  };
  const calModulePath = path => {
    path = path.replace(/\./g, "/");
    path = path.replace(/\/+/g, "/");
    path = path.charAt(0) === "/" ? path.substring(1) : path;
    path =
      path.charAt(path.length - 1) === "/"
        ? path.substring(0, path.length - 1)
        : path;
    return path === "" ? path : "/" + path;
  };
  const moduleName = path => {
    const modulesConfig = [
      "app",
      "customer",
      "enterprise",
      "goods",
      "home",
      "material",
      "order",
      "setting",
      "shop",
      "store"
    ];
    let moduleName = "";
    modulesConfig.forEach(item => {
      if (path.indexOf("/" + item + "/") !== -1) {
        moduleName = item;
      }
    });
    return moduleName;
  };
  const componentInfo = yargs.argv.fileName.split(":"); // 冒号用来确认文件名和路径 冒号前为要创建的文件名，冒号后为此文件的路径 此规则默认为畅移所有项目
  const name = componentInfo[0]; // 默认为文件名
  const parentPath = (componentInfo.length > 1 && componentInfo[1]) || ""; // 判断是否需要生成路径
  const fileName = name.replace(/([A-Z])/g, "-$1").toLowerCase(); // 文件名称转化为小写
  const destPath = path.join(resolveToDemoPages(), parentPath, fileName); // 文件生成的整体路径

  const nameOk = name.replace(/-(\w)/g, function(all, letter) {
    return letter.toUpperCase();
  }); // 文件名称转化为大驼峰
  // console.log(
  //   "文件名称小驼峰",
  //   nameOk,
  //   "文件名默认",
  //   fileName,
  //   "文件名大驼峰",
  //   cap(nameOk),
  //   "指定创建的文件路径",
  //   parentPath,
  //   "指定创建的文件路径加'/'",
  //   calModulePath(parentPath),
  //   "最终文件路径",
  //   destPath,
  //   "模块名称",
  //   moduleName(destPath)
  // );
  return gulp
    .src(paths.blankTemplates)
    .pipe(
      template({
        name: nameOk, //文件名称小驼峰
        fileName: fileName, //文件名默认
        upCaseName: cap(nameOk), // 文件名大驼峰
        parentPath: parentPath,
        modulePath: calModulePath(parentPath),
        moduleName: moduleName(destPath)
      })
    )
    .pipe(
      rename(path => {
        path.basename = path.basename.replace("temp", fileName);
      })
    )
    .pipe(gulp.dest(destPath));
});
