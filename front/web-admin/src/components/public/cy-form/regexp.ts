// 正则方法
const reg = {
  emailReg: /\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*/,
  phoneReg: /^(0|86|17951)?(13[0-9]|15[012356789]|166|17[3678]|18[0-9]|14[57])[0-9]{8}$/,
  idcardReg: /^(^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$)|(^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])((\d{4})|\d{3}[Xx])$)$/,
  numberReg: /^\d+$/,
  urlReg: /^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i,
  zhReg: /[\u4e00-\u9fa5]/gm,
  enReg: /^[a-z]+$/i,
  ennumReg: /^[a-z0-9]+$/i,
  en_numReg: /^\w+$/
};

// validate 接口
interface IObjInterface {
  label: string;
  ruleName: string;
  params: any;
  [propName: string]: any;
}
interface IRulesInterface {
  trigger: string;
  required?: boolean;
  [propName: string]: any;
}
const insideArr = [
  "email",
  "phone",
  "idcard",
  "number",
  "url",
  "zh",
  "en",
  "ennum",
  "en_num",
  "trims"
];

export const regrexFun = function(originRules: Array<any>, innerRules: any) {
  const validate: Array<any> = [];
  // form常规校验
  const normalvalidatorList: Array<IRulesInterface> = [];
  // 自定义校验
  // new Error("required 校验请使用form-item内置校验");
  const validatorList = originRules.filter(item => {
    Object.prototype.hasOwnProperty.call(item, "validator") ||
      normalvalidatorList.push(item);
    item.required &&
      console.error("form => rule required 校验请使用form-item内置校验");
    (item.min || item.max) &&
      console.error("form => rule min | max 校验请使用form-item内置校验");
    return Object.prototype.hasOwnProperty.call(item, "validator");
  });

  // required
  if (innerRules.params["required"]) {
    validate.push({
      message: `请${innerRules.trigger === "blur" ? "输入" : "选择"}${
        innerRules.label
      }`,
      trigger: innerRules.trigger,
      required: true
    });
  }
  // min & max
  if (innerRules.params["min"] || innerRules.params["max"]) {
    const params: IRulesInterface = {
      trigger: "blur"
    };
    innerRules.params["min"] =
      innerRules.params["min"] && Number(innerRules.params["min"]);
    innerRules.params["max"] =
      innerRules.params["max"] && Number(innerRules.params["max"]);
    if (innerRules.params["min"] && innerRules.params["max"]) {
      params.message = `${innerRules.label}长度应在${innerRules.params["min"]}到${innerRules.params["max"]}之间`;
      params.min = innerRules.params["min"];
      params.max = innerRules.params["max"];
    } else if (innerRules.params["min"]) {
      params.message = `${innerRules.label}长度不能小于${innerRules.params["min"]}`;
      params.min = innerRules.params["min"];
    } else if (innerRules.params["max"]) {
      params.message = `${innerRules.label}长度不能大于${innerRules.params["max"]}`;
      params.max = innerRules.params["max"];
    }
    validate.push(params);
  }
  // 常规rule直接push
  validate.push(...normalvalidatorList);
  // email\phone\idcard...
  // 存在才执行
  if (
    insideArr.some(item => {
      return innerRules.params[item];
    })
  ) {
    const regfun = function(_innerRules: any) {
      return function(rule: any, value: any, callback: any) {
        value = value + "";
        // 以下几个常规提示
        if (
          value &&
          !["email", "phone", "idcard", "url"].some(item => {
            if (_innerRules.params[item] && !value.match(reg[item + "Reg"])) {
              callback(new Error(`请输入正确的${_innerRules.label}`));
              return true;
            }
          })
        ) {
          // 剩余特殊提示
          if (_innerRules.params.number && !value.match(reg.numberReg)) {
            callback(new Error(`${_innerRules.label}仅支持输入数字`));
            return;
          }
          if (_innerRules.params.zh && !value.match(reg.zhReg)) {
            callback(new Error(`${_innerRules.label}仅支持输入中文字符`));
            return;
          }
          if (_innerRules.params.en && !value.match(reg.enReg)) {
            callback(new Error(`${_innerRules.label}仅支持输入英文字符`));
            return;
          }
          if (_innerRules.params.ennum && !value.match(reg.ennumReg)) {
            callback(new Error(`${_innerRules.label}仅支持输入数字和英文字符`));
            return;
          }
          if (_innerRules.params.en_num && value.match(reg.en_numReg)) {
            callback(
              new Error(`${_innerRules.label}仅支持输入数字、英文、下划线字符`)
            );
            return;
          }
          if (_innerRules.params.trims && value !== value.trim()) {
            callback(new Error(`请不要输入首尾空格`));
            return;
          }
          validatorList.length || callback();
          validatorList.forEach(item => {
            typeof item.validator === "function" &&
              item.validator(rule, value, callback);
          });
        }
      };
    };
    validate.push({
      validator: regfun(innerRules),
      trigger: "blur"
    });
  } else {
    validate.push(...validatorList);
  }
  if (validate.length === 0) validate.push(...originRules);

  return validate;
};
