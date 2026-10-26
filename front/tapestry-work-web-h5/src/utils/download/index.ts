import axios from 'axios';
import { readFilename } from './util';

/**
 * 下载
 * @param filename 文件名
 * @param content 文件内容
 * @param type 文件类型
 */
export function download(content: string, type = 'text/plain', filename: string): void {
  const blob = new Blob([content], { type });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.style.display = 'none';
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  console.log(link);
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
/**
 * 导出网络文件
 * @param filename 文件名
 * @param url 网络地址
 */
export async function downloadFile(url: string, options?: IDownloadFileOptions) {
  const { filename, onProgress, ...config } = options || {};
  const response = await axios({
    url,
    ...config,
    method: 'GET',
    responseType: 'blob',
    onDownloadProgress: ({ total, loaded }) => {
      const percent = total ? Math.round((loaded * 100) / total) : 0;
      console.log(`下载进度：${percent}%`);
      onProgress && onProgress(percent);
    },
  });
  console.log(response);
  
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

export function getFileData(url: string) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, { responseType: 'blob' })
      .then((res) => {
        const reader = new FileReader();
        reader.onload = () => {
          // 将buffer转换成base64
          resolve(reader.result);
        };
        reader.readAsArrayBuffer(res.data);


      })
      .catch(reject);
  });
}