export const formatDate = function (timestamp: number, format: string) {
  if (!format) {
    format = 'yyyy-MM-dd hh:mm:ss';
  }
  const realDate = new Date(timestamp);

  function timeFormat(num) {
    return num < 10 ? '0' + num : num;
  }
  const date = [
    ['M+', timeFormat(realDate.getMonth() + 1)],
    ['d+', timeFormat(realDate.getDate())],
    ['h+', timeFormat(realDate.getHours())],
    ['m+', timeFormat(realDate.getMinutes())],
    ['s+', timeFormat(realDate.getSeconds())],
    ['q+', Math.floor((realDate.getMonth() + 3) / 3)],
    ['S+', realDate.getMilliseconds()],
  ];
  const reg1 = new RegExp('(y+)', 'i').exec(format);
  // console.log(reg1[0]);
  if (reg1) {
    format = format.replace(
      reg1[1],
      (realDate.getFullYear() + '').substring(4 - reg1[1].length),
    );
  }
  for (let i = 0; i < date.length; i++) {
    const k = date[i][0];
    const v = date[i][1];

    const reg2 = new RegExp('(' + k + ')').exec(format);
    if (reg2) {
      format = format.replace(
        reg2[1],
        reg2[1].length === 1 ? v : ('00' + v).substring(('' + v).length),
      );
    }
  }
  return format;
};