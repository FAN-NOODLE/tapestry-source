// http://eslint.org/docs/user-guide/configuring

module.exports = {
  root: true,
  parser: 'babel-eslint',
  parserOptions: {
    sourceType: 'module'
  },
  // https://github.com/feross/standard/blob/master/RULES.md#javascript-standard-style
  extends: [
    'standard',
    'plugin:promise/recommended',
    'plugin:import/recommended'
  ],
  env: {
    browser: true,
    node: true,
    es6: true
  },
  // globals
  globals: {
    __DEV__: true,
    App: true,
    Page: true,
    wx: true,
    getApp: true,
    Component: true,
    Behavior: true,
    getCurrentPages: true
  },
  settings: {
    'import/resolver': {
      webpack: {
        config: 'webpack.config.babel.js'
      }
    }
  },
  // add your custom rules here
  rules: {
    // always allow new operator
    'no-new': 0, //禁止在使用new构造一个实例后不赋值
    // always add dangling comma for multine
    'comma-dangle': [
      'error',
      {
        arrays: 'never',
        objects: 'never',
        imports: 'never',
        exports: 'never',
        functions: 'ignore'
      }
    ], //对象字面量项尾不能有逗号
    // allow paren-less arrow functions
    'arrow-parens': 0, // 箭头函数用小括号括起来
    // allow async-await
    'generator-star-spacing': 0, // 生成器函数*的前后空格
    // allow debugger during development
    'no-debugger': process.env.NODE_ENV === 'build' ? 2 : 0, // 禁止使用debugger
    // disallow multiple empty lines
    'no-multiple-empty-lines': ['error', { max: 2, maxEOF: 1 }], // 空行最多不能超过2行
    // allow customized promise parameter names
    'promise/param-names': 'off', //
    // from boilerplate
    'import/no-absolute-path': 2,
    'import/no-extraneous-dependencies': 2,
    'import/no-mutable-exports': 2,
    'import/newline-after-import': 1,
    'import/unambiguous': 0,
    'promise/avoid-new': 0,
    'promise/no-callback-in-promise': 0,
    'promise/always-return': 0,
    'promise/no-nesting': 0,
    'space-before-function-paren': 0,
    // enforce the consistent use of either backticks, double, or single quotes
    quotes: ['error', 'single', { allowTemplateLiterals: false }], // 引号类型 ''
    camelcase: [0, { properties: 'never' }] // 不强制使用驼峰
  }
}
