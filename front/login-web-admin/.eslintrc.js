module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ["plugin:vue/essential", "@vue/prettier"],
  rules: {
    "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-new": 0, //禁止在使用new构造一个实例后不赋值
    // "comma-dangle": [0, "always-multiline"], //对象字面量项尾不能有逗号
    'comma-dangle': ['error',  {
      "arrays": "never",
      "objects": "never",
      "imports": "never",
      "exports": "never",
      "functions": "ignore"
    }],
    "arrow-parens": 0, //箭头函数用小括号括起来
    "generator-star-spacing": 0, //生成器函数*的前后空格
    "no-multiple-empty-lines": ["error", { max: 2 }], //空行最多不能超过2行
    "promise/param-names": "off", // allow customized promise parameter names
    "one-var": 0, //连续声明
    "computed-property-even-spacing": [0, "never"], //是否允许计算后的键名什么的
    "no-dupe-args": 2, //函数定义的时候不允许出现重复的参数
    "no-dupe-keys": 2, //对象中不允许出现重复的键
    "no-duplicate-case": 2, //switch语句中不允许出现重复的case标签
    "no-empty": 2, //不允许出现空的代码块
    "no-extra-parens": 0, //不允许出现不必要的圆括号
    "no-irregular-whitespace": 2, //不允许出现不规则的空格
    "no-extra-semi": 2, //不允许出现不必要的分号
    "no-lone-blocks": 2, //不允许不必要的嵌套代码块
    "no-loop-func": 2, //不允许在循环语句中进行函数声明
    "no-multi-spaces": 2, //不允许出现多余的空格
    "no-shadow-restricted-names": 2, //js关键字和保留字不能作为函数名或者变量名
    "brace-style": [2, "1tbs", { allowSingleLine: false }], //大括号风格
    "consistent-this": [0, "self"] //当获取当前环境的this是用一样的风格
  },
  parserOptions: {
    parser: "babel-eslint"
  }
};
