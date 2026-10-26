/* eslint-disable import/no-extraneous-dependencies */
import { resolve } from 'path'
import { DefinePlugin, EnvironmentPlugin, optimize } from 'webpack'
import WXAppWebpackPlugin from 'wxapp-webpack-plugin'
import CopyWebpackPlugin from 'copy-webpack-plugin'
import CleanWebpackPlugin from 'clean-webpack-plugin'
import ENV_DEV from './deploy/dev/env.dev.json'
import ENV_PROD from './deploy/prod/env.prod.json'


// 定义环境数据
let config = {
  dev: ENV_DEV,
  build: ENV_PROD
}


const relativeFileLoader = (ext = '[ext]') => [
  {
    loader: 'file-loader',
    options: {
      publicPath: '',
      useRelativePath: true,
      name: `[name].${ext}`,
      emitFile: false
    }
  },
  {
    loader: 'file-loader',
    options: {
      publicPath: '',
      context: resolve('src'),
      name: `[path][name].${ext}`
    }
  }
]

export default (env, argv) => {
  const { outputPath } = argv
  return {
    entry: {
      app: './src/app.js'
    },
    output: {
      filename: '[name].js',
      publicPath: '/',
      path: resolve(outputPath || 'dist')
    },
    module: {
      rules: [
        {
          test: /\.js$/,
          include: /src/,
          exclude: /lib/,
          use: [
            {
              loader: 'babel-loader',
              options: {
                presets: ['env', 'stage-2'],
                plugins: [
                  [
                    'transform-runtime',
                    {
                      polyfill: true,
                      regenerator: true
                    }
                  ],
                  [
                    'transform-decorators-legacy',
                    {
                      legacy: true
                    }
                  ]
                ]
              }
            }
            // 'eslint-loader'
          ]
        },
        {
          test: /\.wxs$/,
          include: /src/,
          use: [...relativeFileLoader(), 'babel-loader', 'eslint-loader']
        },
        {
          test: /\.(scss|sass|wxss)$/,
          include: /src/,
          use: [
            ...relativeFileLoader('wxss'),
            'extract-loader',
            'css-loader',
            'sass-loader'
          ]
        },
        {
          test: /\.(png|jpe?g|gif|svg|eot|ttf|woff|woff2)$/,
          include: /src/,
          loader: 'url-loader',
          options: {
            limit: 10000,
            fallback: 'file-loader',
            name: './images/[name].[ext]'
          }
        },
        {
          test: /\.json$/,
          include: /src/,
          exclude: /ignore/,
          use: relativeFileLoader()
        },
        {
          test: /\.wxml$/,
          include: resolve('src'),
          use: [
            ...relativeFileLoader(),
            {
              loader: 'wxml-loader',
              options: {
                root: resolve('src'),
                minimize: false
              }
            }
          ]
        }
      ]
    },
    plugins: [
      new EnvironmentPlugin({
        NODE_ENV: 'development'
      }),
      new DefinePlugin({
        'process.env': JSON.stringify(config[env])
      }),
      new CleanWebpackPlugin([
        `./${outputPath || 'dist'}/*`
      ]),
      new WXAppWebpackPlugin({
        clear: false
      }),
      new optimize.ModuleConcatenationPlugin(),
      new CopyWebpackPlugin(['src/project.config.json'])
    ].filter(Boolean),
    // devtool: isDev ? 'source-map' : false,
    resolve: {
      extensions: ['.js', '.wxml', '.sass', '.scss'],
      modules: [resolve(__dirname, 'src'), 'node_modules'],
      alias: {
        '@': resolve('src')
      }
    },
    watchOptions: {
      ignored: /dist|build/,
      aggregateTimeout: 300
    }
  }
}
