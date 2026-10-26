import Axios from "axios";

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
export function getFilePath(file: File|string, folder:string):string {
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

// export function getVideoDuration(file: File):Promise<number> {
//   return new Promise((resolve, reject) => {
//     const video = document.createElement('video');
//     video.preload = 'metadata';
//     video.src = URL.createObjectURL(file);
//     video.onloadedmetadata = () => {
//       resolve(parseInt(video.duration.toString()));
//     };
//     video.onerror = () => {
//       reject(new Error('Failed toload video metadata'));
//     }
//   })
// }

export function getFileData(url: string):Promise<any> {
  return new Promise((resolve, reject) => {
    Axios
      .get(url, { responseType: 'blob' })
      .then((res) => {
        const reader = new FileReader();
        reader.onload = () => {
          // 将buffer转换成base64
          resolve(reader.result as any);
        };
        reader.readAsArrayBuffer(res.data);


      })
      .catch(reject);
  });
}

/**
 * 浏览器端计算文件 Hash
 * @param file File 对象
 * @param algorithm 算法：SHA-1 / SHA-256 / SHA-512
 * @returns hash 十六进制字符串
 */
export async function getFileHash(
  file: File,
  algorithm: 'SHA-1' | 'SHA-256' | 'SHA-512' = 'SHA-256'
): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsArrayBuffer(file);

    reader.onload = async () => {
      try {
        const arrayBuffer = reader.result as ArrayBuffer;
        // 计算摘要
        const hashBuffer = await crypto.subtle.digest(algorithm, arrayBuffer);
        // 转十六进制
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        resolve(hashHex);
      } catch (err) {
        reject(err);
      }
    };

    reader.onerror = () => reject(reader.error);
  });
}

/**
 * 利用COS CI 媒体快照生成视频封面地址
 * @param videoUrl COS原视频地址
 * @param time 截取第几秒，默认1.5秒
 * @returns 封面图片完整地址
 */
export function getCosVideoCoverUrl(videoUrl: string, time = 1.5): string {
  // 已带参数先截断
  const baseUrl = videoUrl.split('?')[0];
  return `${baseUrl}?ci-process=snapshot&time=${time}`;
}

/**
 * 获取COS视频封面，返回Blob
 */
export async function getCosCoverBlob(videoUrl: string, time = 1.5) {
  const coverUrl = getCosVideoCoverUrl(videoUrl, time);
  const res = await Axios.get(coverUrl, {
    responseType: 'blob'
  });
  return res.data as Blob;
}

/**
 * 获取COS视频封面Base64
 */
export async function getCosCoverBase64(videoUrl: string, time = 1.5) {
  const blob = await getCosCoverBlob(videoUrl, time);
  return new Promise<string>((resolve) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.readAsDataURL(blob);
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
  const { url, quality = 80, maxWidth, maxHeight, format = 'jpg' } = opt;
  const base = url.split('?')[0];
  const params = [`imageMogr2/quality/${quality}`, `format/${format}`];

  if (maxWidth && !maxHeight) {
    params.push(`thumbnail/${maxWidth}x`);
  } else if (maxWidth && maxHeight) {
    params.push(`thumbnail/${maxWidth}x${maxHeight}`);
  } else if(maxHeight && !maxWidth) {
    params.push(`thumbnail/x${maxHeight}`);
  }

  return `${base}?${params.join('/')}`;
}

/**
 * 本地获取视频时长（上传前校验）
 * @param {File} file - input[type=file] 选中的 file 对象
 * @returns {Promise<number>} 秒数
 */
export function getVideoDuration(file: File):Promise<number> {
  return new Promise((resolve, reject) => {
    // 非视频直接拒绝
    if (!file.type.startsWith('video/')) {
      return reject(new Error('不是视频文件'))
    }

    const video = document.createElement('video')
    video.preload = 'metadata'

    // 解决部分浏览器跨域问题
    video.crossOrigin = 'anonymous'

    // 本地文件生成临时 URL
    const url = URL.createObjectURL(file)
    video.src = url

    // 元数据加载完成 → 获取时长
    video.onloadedmetadata = () => {
      const duration = video.duration
      // 释放内存
      URL.revokeObjectURL(url)
      video.remove()
      resolve(duration)
    }

    video.onerror = (err) => {
      URL.revokeObjectURL(url)
      reject(new Error('视频无法读取'))
    }
  })
}