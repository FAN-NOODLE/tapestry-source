// 下载文件到本地的方法
export function downLoad(url: string, name: string) {
  const x = new XMLHttpRequest();
  const newUrl = url + (url.indexOf('?') > -1 ? '&' : '?') + 'timestamp=' + new Date().getTime();
  console.log("download", newUrl, name);
  x.open("GET", newUrl, true);
  x.responseType = "blob";
  x.onload = function() {
    const urlResponse = window.URL.createObjectURL(x.response);
    const a = document.createElement("a");
    a.href = urlResponse;
    a.download = name;
    a.click();
  };
  x.send();
}

// 下载文件到本地的方法
export function downLoadSync(url: string, name: string) {
  return new Promise((resolve, reject) => {
    const x = new XMLHttpRequest();
    const newUrl = url + (url.indexOf('?') > -1 ? '&' : '?') + 'timestamp=' + new Date().getTime();
    console.log("download", newUrl, name);
    x.open("GET", newUrl, true);
    x.responseType = "blob";
    x.onload = function() {
      const urlResponse = window.URL.createObjectURL(x.response);
      const a = document.createElement("a");
      a.href = urlResponse;
      a.download = name;
      a.click();
    };
    x.onerror = function() {
      reject(false)
    };
    x.send();
  });
}



