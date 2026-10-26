// const UglifyJsPlugin = require("uglifyjs-webpack-plugin");
const ANALYZER = false;
const path = require("path");
const TerserPlugin = require("terser-webpack-plugin");
const assetsDir = "static";
// todo: 增加gzip
const CompressionWebpackPlugin = require("compression-webpack-plugin");
const productionGzipExtensions = /\.(js|css|json|txt|html|ico|svg)(\?.*)?$/i;
let configPath = () => {
  let isNewDevops = process.env.VUE_APP_DEVOPS === "new";
  const publicPath = isNewDevops ? "/login/" : "/login/";
  const assetsDir = "static";
  // 新的ci-cd流程 outputDir 为 “ci-cd/build/dist/" 旧的ci-cd流程 outputDir 为 deploy/dist/+publicPath
  const outputDir = (isNewDevops ? "ci-cd/build" : "deploy") + "/dist";
  return {
    publicPath: publicPath,
    assetsDir: assetsDir,
    outputDir: outputDir
  };
};
function getAssetPath(assetsDir, filePath) {
  return assetsDir ? path.posix.join(assetsDir, filePath) : filePath;
}
module.exports = {
  devServer: {
    open: true,
    port: 20000
  },
  publicPath: configPath().publicPath,
  assetsDir: configPath().assetsDir, // 放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录
  outputDir: configPath().outputDir,
  filenameHashing: false, // 生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存
  productionSourceMap: false, // 正式环境不生成 source map
  chainWebpack: config => {
    if (ANALYZER) {
      // 打包分析模式
      config
        .plugin("webpack-bundle-analyzer")
        .use(require("webpack-bundle-analyzer").BundleAnalyzerPlugin)
        .end();
    }
    const filename = getAssetPath(assetsDir, `/js/[name]-[hash].js`);
    config
      .mode(process.env.NODE_ENV)
      .devtool(false)
      .output.filename(filename)
      .chunkFilename(filename);
    config.plugin("extract-css").tap(() => [
      {
        filename: getAssetPath(assetsDir, `/css/[name]-[hash].css`),
        chunkFilename: getAssetPath(assetsDir, `/css/[name]-[hash].css`)
      }
    ]);
  },
  configureWebpack: config => {
    if (process.env.NODE_ENV === "production") {
      // 为生产环境修改配置...
      config.plugins = [
        ...config.plugins,
        new TerserPlugin({
          // 压缩Javascript
          cache: true, // ？？
          parallel: true, // 使用多进程并发运行以提高构建速度
          sourceMap: false, // ？？
          terserOptions: {
            compress: {
              drop_debugger: true, // true:删除debugger 语句
              drop_console: true // true:放弃对console.*函数的调用
            }
          }
        }),
        // todo:配置compression-webpack-plugin压缩
        new CompressionWebpackPlugin({
          filename: "[path][base].gz",
          algorithm: "gzip",
          // 要压缩的文件（正则）
          test: productionGzipExtensions,
          // 最小文件开启压缩
          threshold: 10240,
          minRatio: 0.8
        })
      ];
    }
  }
};
