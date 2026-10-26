const arithmetic = {
  baseA: {
    // 解密
    decode(code) {
      code = unescape(code);
      let c = String.fromCharCode(code.charCodeAt(0) - code.length);
      for (let i = 1; i < code.length; i++) {
        c += String.fromCharCode(code.charCodeAt(i) - c.charCodeAt(i - 1));
      }
      return c;
    },
    // 加密
    encode(code) {
      let c = String.fromCharCode(code.charCodeAt(0) + code.length);
      for (let i = 1; i < code.length; i++) {
        c += String.fromCharCode(code.charCodeAt(i) + code.charCodeAt(i - 1));
      }
      return `baseA:${escape(c)}`;
    }
  }
};

// storage 存储敏感数据都要使用密文
export default func => {
  if (func === "encode") {
    return (type, str) => arithmetic[type].encode(str);
  }
  if (func === "decode") {
    return str => arithmetic[str.split(":")[0]].decode(str.split(":")[1]);
  }
};
