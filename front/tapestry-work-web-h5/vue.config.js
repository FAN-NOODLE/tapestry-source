const path = require("path");
const openInEditor = require('launch-editor-middleware');
const CompressionPlugin = require('compression-webpack-plugin')
const TerserPlugin = require("terser-webpack-plugin");
let configPath = () => {
  // 新的ci-cd流程 publicPath 为 “/work/” 旧的ci-cd流程 publicPath 为 “/work/+tag”
  const publicPath = "/work/";
  // 新的ci-cd流程 assetsDir 为 “版本号+static” 旧的ci-cd流程 assetsDir 为 “static”
  // const assetsDir = `${process.env.VUE_APP_VERSION}/static`;
  const assetsDir = `static`;
  // 新的ci-cd流程 outputDir 为 “ci-cd/build/dist/+publicPath(新的publicPath)" 旧的ci-cd流程 outputDir 为 deploy/dist/+publicPath(旧的publicPath)
  // const outputDir = "ci-cd/build/dist" + publicPath;
  const outputDir = "deploy/dist";

  // 环境动态domain配置
  const envConfig = {
    test: "-test",
    dev: "-dev",
    pro: "",
    uat: "",
  };
  return {
    publicPath: publicPath,
    assetsDir: assetsDir,
    outputDir: outputDir,
  };
};
module.exports = {
  publicPath: configPath().publicPath,
  assetsDir: configPath().assetsDir, // 放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录
  outputDir: configPath().outputDir,
  filenameHashing: true, // 生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存
  productionSourceMap: process.env.NODE_ENV !== "production", // 正式环境不需要
  lintOnSave: false, // 关闭eslint规范
  css: {
    loaderOptions: {
      less: {
        // 若 less-loader 版本小于 6.0，请移除 lessOptions 这一级，直接配置选项。
        lessOptions: {
          modifyVars: {
            // 直接覆盖变量
            "tabs-default-color": " #0095f6",
            "dropdown-menu-height": "32px",
            "dropdown-menu-background-color": "#f7f8fa",
            "dropdown-menu-box-shadow": "none",
            "nav-bar-title-font-size": "13px",
            // 或者可以通过 less 文件覆盖（文件路径为绝对路径）
            // hack: `true; @import "theme.less";`,
          },
        },
      },
      postcss: {
        plugins: [
          require("autoprefixer"),
          require("postcss-pxtorem")({
            rootValue: 37.5, // 换算的基数
            selectorBlackList: ["layout"], // 忽略转换正则匹配项
            propList: ["*", "!border"],
            unitPrecision: 5, //保留小数位
            minPixelValue: 2, //转换的最小单位
          }),
        ],
      },
    },
  },
  parallel: false, //防止vant打包后样式丢失
  devServer: {
    before (app) {
      app.use('/__open-in-editor', openInEditor())
    },
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Headers': 'EP-Version,Enterprise-Hash'
    },
    port: 8080,
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_ROOT,
        changeOrigin: true, //是否跨域
        ws: true, 
        pathRewrite: {
         '/api': process.env.VUE_APP_API_ROOT
        }
      },
      '/robot': {
        target: 'https://qyapi.weixin.qq.com',
        changeOrigin: true, //是否跨域
        ws: true,
        pathRewrite: {
          '/robot': ''
        }
      }
    }
  },
  chainWebpack: (config) => {
    //删除懒加载模块的 prefetch preload，降低带宽压力
    config.plugins.delete("prefetch").delete("preload");
    config.resolve.alias.set("@", path.resolve("src"));
    // 配置svg-icon loader
    config.module.rule("svg").exclude.add(path.resolve("@/icons")).end();
    config.module
      .rule("icon")
      .test(/\.svg$/)
      .include.add(path.resolve("@/icons"))
      .end()
      .use("svg-sprite-loader")
      .loader("svg-sprite-loader")
      .options({
        symbolId: "icon-[name]",
      })
      .end();
  },
  // todo:gzip压缩
  configureWebpack: (config) => {
    if (process.env.NODE_ENV === 'production') {
      config.plugins.push(
        new TerserPlugin({
          // 压缩Javascript
          cache: true, // ？？
          parallel: true, // 使用多进程并发运行以提高构建速度
          sourceMap: false, // ？？
          terserOptions: {
            compress: {
              drop_debugger: true, // true:删除debugger 语句
              drop_console: false // true:放弃对console.*函数的调用
            }
          }
        }),
        new CompressionPlugin({
          // gzip压缩配置
          test: /\.js$|\.html$|\.css/, // 匹配文件名
          threshold: 10240, // 对超过10kb的数据进行压缩
          deleteOriginalAssets: false, // 是否删除原文件
        })
      )
      // 解决css文件打包顺序问题
      const miniCssExtractPlugin = config.plugins.find(
        plugin => plugin.constructor.name === 'MiniCssExtractPlugin'
      );
      if (miniCssExtractPlugin) {
        miniCssExtractPlugin.options.ignoreOrder = true;
      }
    }
  }
};
