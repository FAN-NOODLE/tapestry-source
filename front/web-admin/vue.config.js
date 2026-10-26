const path = require("path");
const assetsDir = "static";
// const OptimizeCSSPlugin = require("optimize-css-assets-webpack-plugin");
const isProduction = process.env.NODE_ENV === "production";
// todo:引入cdn
const webpackConfig = require("./webpack-config.js")
// todo: 编译缓存
const HardSourceWebpackPlugin = require('hard-source-webpack-plugin');

function resolve(dir) {
  return path.join(__dirname, dir);
}
function getAssetPath(assetsDir, filePath) {
  return assetsDir ? path.posix.join(assetsDir, filePath) : filePath;
}

const configPath = () => {
  const publicPath = "/web-admin";
  const assetsDir = "static";
  const outputDir = "dist";
  return {
    publicPath: publicPath,
    assetsDir: assetsDir,
    outputDir: outputDir,
  };
};

module.exports = {
  lintOnSave: false,
  devServer: {
    open: true,
    port: 20000,
    // 当本地需要调试生产环境时proxy代理解决接口跨域
    proxy: {
      "/crs-prod-api": {
        target: "https://console-api-crs.vchangyi.com/",
        secure: true,
        changeOrigin: true,
        pathRewrite: {
          "^/crs-prod-api": "",
        },
      },
    },
  },
  publicPath: configPath().publicPath,
  // 放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目 录
  assetsDir: configPath().assetsDir,
  outputDir: configPath().outputDir,
  // 生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存
  filenameHashing: false,
  productionSourceMap: false,
  css: {
    extract: true,
    sourceMap: !isProduction,
    requireModuleExtension: true,
  },

  chainWebpack: (config) => {
    // 移除 prefetch 插件
    config.plugins.delete("prefetch");
    // 配置svg-icon loader
    config.module
      .rule("svg")
      .exclude.add(resolve("src/icon"))
      .end();
    config.module
      .rule("icon")
      .test(/\.svg$/)
      .include.add(resolve("src/icon"))
      .end()
      .use("svg-sprite-loader")
      .loader("svg-sprite-loader")
      .options({
        symbolId: "icon-[name]",
      })
      .end();
    // 配置打包文件名增加hash值
    const filename = getAssetPath(assetsDir, `/js/[name]-[hash].js`);
    config
      .mode(process.env.NODE_ENV)
      .devtool(false)
      .output.filename(filename)
      .chunkFilename(filename);
    // css打包文件增加hash值
    config.plugin("extract-css").tap(() => [
      {
        filename: getAssetPath(assetsDir, `/css/[name]-[hash].css`),
        chunkFilename: getAssetPath(assetsDir, `/css/[name]-[hash].css`),
      },
    ]);
    // 注入cdn star
    config.plugin('html').tap(args => {
      // 生产环境或本地需要cdn时，才注入cdn
      args[0].cdn = webpackConfig.cdn
      return args
    })
  },
  configureWebpack: (config) => {
    if(isProduction){
      config.plugins = [
        ...config.plugins,
        ...webpackConfig.webpackPlugins,
        new HardSourceWebpackPlugin()
      ];
    }else{
      config.plugins = [
        ...config.plugins
      ];
    }
    // 解决css文件打包顺序问题
    const miniCssExtractPlugin = config.plugins.find((plugin) => plugin.constructor.name === "MiniCssExtractPlugin");
    if (miniCssExtractPlugin) {
      miniCssExtractPlugin.options.ignoreOrder = true;
    }
    // 配置cdn
    config.externals = webpackConfig.cdn.externals
  },
  pluginOptions: {
    "style-resources-loader": {
      preProcessor: "scss",
      patterns: [
        // 注入的资源的路径
        path.resolve(__dirname, "src/style/color/index.scss"),
        path.resolve(__dirname, "src/style/mixins/index.scss"),
      ],
    },
  },
};
