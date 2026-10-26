// 格式化文件大小数据
export default {
  key: "formatFileSize",
  func: (size: number): string => {
    if (size < 1024) {
      return `${size}B`;
    }
    if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(2)}KB`;
    }
    if (size < 1024 * 1024 * 1024) {
      return `${(size / (1024 * 1024)).toFixed(2)}MB`;
    }
    return `${(size / (1024 * 1024 * 1024)).toFixed(2)}GB`;
  }
};
