import axios from 'axios';
import { readFilename } from './util';
import { VUE_APP_API_ROOT_RE } from "@/utils/env-config"
/**
 * 下载
 * @param filename 文件名
 * @param content 文件内容
 * @param type 文件类型
 */
export function download(content: string | Blob, type = 'text/plain', filename?: string): void {
  const blob = typeof content ==='string'? new Blob([content], { type }) : content;
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.style.display = 'none';
  link.href = url;
  if(filename) {
    link.download = filename;
  }
  document.body.appendChild(link);
  link.click();
  // 延时确保不会中断下载
  setTimeout(() => {
    // 创建完毕后，移除 <a> 标签，并释放用完的 URL
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }, 0);
}

/**
 * 导出json文件
 * @param filename 文件名
 * @param content 文件内容
 */
export function downloadJson(content: any, filename: string): void {
  download(JSON.stringify(content, null, 2), 'application/json', filename);
}

export interface IDownloadFileOptions {
  filename?: string;
  onProgress?: (progress: number) => void;
  [key: string]: any;
}

export async function getFileBlob(url: string, options?: IDownloadFileOptions) {
  const { onProgress, ...config } = options || {};
  const response = await axios({
    url: url + (url.indexOf('?') > -1 ? '&' : '?') + 'timestamp=' + new Date().getTime(),
    responseType: 'blob',
    ...config,
    onDownloadProgress: ({ total, loaded }) => {
      const percent = total ? Math.round((loaded * 100) / total) : 0;
      console.log(`下载进度：${percent}%`);
      onProgress && onProgress(percent);
    },
  });

  return response
}
/**
 * 导出网络文件
 * @param filename 文件名
 * @param url 网络地址
 */
export async function downloadFile(url: string, options?: IDownloadFileOptions) {
  const { filename } = options || {};
  const response = await getFileBlob(url, options)
  
  download(
    response.data,
    response.headers['content-type'],
    filename || readFilename(response.headers['content-disposition']),
  );
  return response;
}



export function getFileBase64(url: string) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, { responseType: 'blob' })
      .then((res) => {
        const reader = new FileReader();
        reader.onload = () => {
          const base64Data = (reader.result as string).split(',')[1];
          resolve(base64Data);
        };
        reader.readAsDataURL(res.data);
      })
      .catch(reject);
  });
}


export function getBlobUrl(url: string): Promise<string> {
  const encoder = new TextEncoder();
  const encodedData = encoder.encode(url);
  // 使用 btoa 编码字节数组为 Base64 字符串
  const encodedText = btoa(String.fromCharCode(...encodedData));
  // 获取随机2位字符串
  const randomStr = Math.random().toString(36).substring(2, 4);
  
  const data = `${VUE_APP_API_ROOT_RE}common/get-file-by-url?url=${randomStr}${encodedText}`;
    
  return new Promise((resolve, reject) => {
    const x = new XMLHttpRequest();
    x.open('GET', data, true);
    x.responseType = 'blob';
    x.onload = function () {
      const src = window.URL.createObjectURL(x.response)
      resolve(src);
    };
    x.onerror = reject;
    x.send();
    
  });
}