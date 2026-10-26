const BaseData = {
  // 基础数据集合
  nil: {
    // 不跳转
    value: 0
  },
  goods: {
    // 商品
    value: 1
  },
  class: {
    // 分类
    value: 2
  },
  local_coupon: {
    // 站内优惠券
    value: 3
  },
  weixin_coupon: {
    // 微信优惠券
    value: 4
  },
  reduce_coupon: {
    // 立减金
    value: 5
  },
  goods_sku: {
    // 直接加购
    value: 6
  },
  activity: {
    // 活动页
    value: 7
  },
  coupons: {
    // 优惠券组
    value: 8
  },
  FlashSale: {
    // 限时购
    value: 9
  },
  Newcomer: {
    // 新人购
    value: 10
  }
};

// 检测是否有重复 value
let repeatCheck = () => {
  let keys = new Set(),
    values = new Set();
  for (let k in BaseData) {
    keys.add(k);
    values.add(BaseData[k]);
  }
  return keys.length === values.length;
};

export default {
  // 转换器，将输入的 key 转化为 value
  convert(params) {
    if (repeatCheck()) {
      if (params.constructor === Array) {
        return params.map(key => BaseData[key].value);
      } else if (params.constructor === String) {
        return BaseData[params].value;
      } else if (params.constructor === Number) {
        let res = "";
        for (let k in BaseData) {
          if (BaseData[k].value === params) {
            res = k;
          }
        }
        return res;
      } else {
        throw new Error("错误：Params must be Array or String!");
      }
    } else {
      throw new Error("错误：Assemble has repeat value!");
    }
  }
};
