import { getVideoPoster } from "../tools";

function toFixed(num:number):string {
  return num > 10 ? num.toString() : "0" + num
}
/**
 * 获取文件目录
 * 如果传入的是文件则格式化文件名称
 * 否则直接使用传入的字符串作为文件名称
 * @param file 
 * @param folder 
 * @returns 
 */
export function getFilePath(file: File|string, folder = 'uploads'):string {
  const isFile = typeof file !== 'string';
  const chars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678";
  const date = new Date();
  const year = date.getFullYear();
  const month = toFixed(date.getMonth() + 1);
  const day = toFixed(date.getDate());
  const paths = [folder||'uploads', year, month, day];
  if(isFile) {
    let filename = "";
    for (let i = 0; i < 32; i++) {
      filename += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    filename += `.${file.name.split(".").pop()}`
    paths.push(filename);
  } else {
    paths.push(file);
  }
  return paths.join('/')
}

export function getVideoDuration(file: File):Promise<number> {
  return new Promise((resolve, reject) => {
    const video = document.createElement('video');
    video.preload = 'metadata';
    video.src = URL.createObjectURL(file);
    video.onloadedmetadata = () => {
      resolve(parseInt(video.duration.toString()));
    };
    video.onerror = () => {
      reject(new Error('Failed toload video metadata'));
    }
  })
}

export function getVideoCover(file: File, url): Promise<string> {
  return new Promise((resolve, reject) => {
    const video = document.createElement('video');
    video.src = url || URL.createObjectURL(file);
    getVideoPoster(video).then((res:any) => {
      console.log(res);
      
      resolve(res.base64);
    }).catch((error) => {
      reject(new Error('Failed to get video poster: ' + error.message));
    })
  });
}
/**
 * COS 图片压缩 URL 生成
 * @param url 原COS图片地址
 * @param quality 画质 1-100 默认70
 * @param maxWidth 限制最大宽度
 * @param format 输出格式 jpg/png/webp
 */
export function getCosCompressUrl(
  opt: {
    url: string;
    quality?: number;
    maxWidth?: number;
    maxHeight?: number;
    format?: 'jpg' | 'png' | 'webp';
  }
): string {
  const { url, quality = 60, maxWidth, maxHeight, format = 'jpg' } = opt;
  const base = url.split('?')[0];
  const params = [`imageMogr2/quality/${quality}`, `format/${format}`];

  if (maxWidth && !maxHeight) {
    params.push(`thumbnail/${maxWidth}x`);
  } else if (maxWidth && maxHeight) {
    params.push(`thumbnail/${maxWidth}x${maxHeight}`);
  } else if(maxHeight && !maxWidth) {
    params.push(`thumbnail/x${maxHeight}`);
  }
  console.log(`${base}?${params.join('/')}`);
  
  return `${base}?${params.join('/')}`;
}