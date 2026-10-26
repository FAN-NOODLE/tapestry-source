// 格式化文本，显示指定长度的文本，超过显示...
export default {
  key: "formatText",
  func: (text: string, length: number): string => {
    if (text.length > length) {
      return `${text.substr(0, length)}...`;
    }
    return text;
  }
};
