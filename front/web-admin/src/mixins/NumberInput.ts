import { Vue, Component } from "vue-property-decorator";
@Component
export default class NumberInput extends Vue {
  /**
   * 校验整数 input事件
   * @description 情况一：深度改变对象中的数字
   * @author LeeYunxiang
   * @param {Object} object 数字的父对象
   * @param {String} target 要监听的字段
   * @param {Number} value 数字
   * @description 情况二：改变的仅是data中的数字
   * @param {String} target 要监听的字段
   * @param {Number} value 数字
   * */
  public _checkInt(...arg: Array<any>): void {
    const reg = /[^0-9$]/g;
    this.$nextTick(() => {
      if (arg.length === 3) {
        const object = arg[0];
        const target = arg[1];
        const value = arg[2];
        let number = value.replace(reg, "");
        if (number.split("").includes("$")) {
          number = number.substr(0, number.length - 1);
        }
        this.$set(object, target, number);
      } else {
        const target = arg[0];
        const value = arg[1];
        let number = value.replace(reg, "");
        console.log(number);
        if (number.split("").includes("$")) {
          number = number.substr(0, number.length - 1);
        }
        this.$data[target] = number;
        console.log(this.$data, target);
      }
    });
  }
  /**
   * 校验正整数 input事件
   * @description 情况一：深度改变对象中的数字
   * @author LeeYunxiang
   * @param {Object} object 数字的父对象
   * @param {String} target 要监听的字段
   * @param {Number} value 数字
   * @description 情况二：改变的仅是data中的数字
   * @param {String} target 要监听的字段
   * @param {Number} value 数字
   * */

  public _checkPositive(...arg: Array<any>): void {
    this.$nextTick(() => {
      if (arg.length === 3) {
        const object = arg[0];
        const target = arg[1];
        const value = arg[2];
        const flag = new RegExp("^[1-9][0-9]{0,}$").test(value);
        flag || this.$set(object, target, "");
      } else {
        const target = arg[0];
        const value = arg[1];
        const flag = new RegExp("^[1-9][0-9]{0,}$").test(value);
        if (!flag) {
          this.$data[target] = "";
        }
      }
    });
  }

  /**
   * 限制数字输入小数点后几位 input事件
   * @description 情况一：深度改变对象中的数字
   * @author ChenYaohong
   * @param {Number} limit 允许输入小数点后几位 例：1代表允许输入小数点后1位
   * @param {Number} nums input框输入的值
   * @param {Number} object 数字的父对象
   * @param {String} target 要监听的字段
   * @description 情况二：改变的仅是data中的数字
   * @param {Number} limit 允许输入小数点后几位 例：1代表允许输入小数点后1位
   * @param {Number} nums input框输入的值
   * @param {String} target 要监听的字段
   * */
  public _checkFloat(...arg: Array<any>): void {
    const onlyDecimal = (num: string, limit: number) => {
      let number = null;
      if (num.indexOf(".") != -1) {
        const str_ = num.substr(num.indexOf(".") + 1);
        if (str_.indexOf(".") != -1) {
          number = num.substr(0, num.indexOf(".") + str_.indexOf(".") + 1);
        }
      }
      const array = num.split(".");
      if (!!array[1] && array[1].length > limit) {
        array[1] = array[1].substr(0, limit);
        number = array[0] + "." + array[1];
      }
      return number;
    };
    this.$nextTick(() => {
      const limit = arg[0];
      let nums = arg[1];
      if (nums) {
        nums = nums.replace(/[^0-9.]/g, "");
        if (onlyDecimal(nums, limit) != null) {
          nums = onlyDecimal(nums, limit);
        }
      }
      if (nums.split("").includes("$")) {
        nums = nums.substr(0, nums.length - 1);
      }
      if (arg.length === 4) {
        const object = arg[2];
        const target = arg[3];
        this.$set(object, target, nums);
      } else {
        const target = arg[2];
        this.$data[target] = nums;
      }
    });
  }
  /**
   * @description 计算小数点后面有几位数
   * @author LeeYunxiang
   * @param {Number} n 要验证的数字
   *
   */
  _countFloat(n: number) {
    try {
      return n.toString().split(".")[1].length;
    } catch (err) {
      return 0;
    }
  }
}
