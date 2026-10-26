// 压缩文件
const TerserPlugin = require("terser-webpack-plugin");
// todo: 增加gzip
const CompressionWebpackPlugin = require("compression-webpack-plugin");
const productionGzipExtensions = /\.(js|css|json|txt|html|ico|svg)(\?.*)?$/i;
// todo: 版本更新检测
const { WebUpdateNotificationPlugin } = require('@plugin-web-update-notification/webpack')

const cdn = {
  // cdn：模块名称和模块作用域命名（对应window里面挂载的变量名称）
  externals: {
    // vue: "Vue",
    // echarts: "echarts",
    // "js-md5": "md5",
    // wangeditor: "wangEditor",
    // xlsx: "XLSX",
  },
  // cdn的css链接
  css: [],
  // cdn的js链接
  js: [
    // "https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/web-admin/js/vue.min.js",
    // "https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/web-admin/js/echarts.min.js",
    // "https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/web-admin/js/md5.min.js",
    // "https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/web-admin/js/wangEditor.min.js",
    // "https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/web-admin/js/xlsx.full.min.js",
  ],
};
const webpackPlugins = [
  new TerserPlugin({
    cache: true,
    parallel: true,
    sourceMap: false,
    terserOptions: {
      compress: {
        drop_debugger: true,
        drop_console: true,
      },
    },
  }),
  // todo:配置compression-webpack-plugin压缩
  new CompressionWebpackPlugin({
    filename: "[path][base].gz",
    algorithm: "gzip",
    // 要压缩的文件（正则）
    test: productionGzipExtensions,
    // 最小文件开启压缩
    threshold: 10240,
    minRatio: 0.8,
  }),
  // todo: 版本更新检测
  new WebUpdateNotificationPlugin({
    logVersion: false,
    checkInterval:0,
    versionType:'build_timestamp',
    hiddenDismissButton: true,
    notificationConfig:{
      primaryColor:"#000000",
    }
  }),
];
module.exports = {
  cdn: cdn,
  webpackPlugins: webpackPlugins
};
