/**
 * 获取文件名称
 * @param contentDisposition
 * @returns
 */
export function readFilename(contentDisposition: string) {
  let filename = '';
  if (contentDisposition) {
    const name1 = contentDisposition.match(/filename=([^;]*)/)?.[1]; // 获取filename的值
    const name2 = contentDisposition.match(/filename\*=([^;]*)/)?.[1]; // 获取filename*的值

    if (name1) {
      filename = decodeURIComponent(name1);
    }

    if (name2) {
      // 解析 filename* 的 UTF-8''syntax (RFC 6266)
      filename = decodeURIComponent(name2.split("'")[2]);
    }
  }

  return filename;
}

function getMimeTypeFromExtension(extension: string): string {
  const extensionMap: Record<string, string> = {
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    png: 'image/png',
    gif: 'image/gif',
    webp: 'image/webp',
    // 添加其他可能的图片格式和对应的 MIME 类型
  };

  // 转换为小写，以便对比
  const lowercasedExtension = extension.toLowerCase();

  // 返回对应的 MIME 类型，如果没有找到则默认为 'application/octet-stream'
  return extensionMap[lowercasedExtension] || 'application/octet-stream';
}

/**
 * 异步加载图片并返回HTMLImageElement对象
 * @param {string} url - 图片的URL
 * @returns {Promise<HTMLImageElement>} - 加载完成的图片对象的Promise
 */
export async function loadImage(url: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.onload = function () {
      resolve(img);
    };
    img.onerror = function (error) {
      reject(error);
    };
    img.src = url.includes('?') ? url : url + '?v=' + Date.now();
  });
}
interface IImageData {
  url: string;
  blob: Blob | null;
  base64: string;
}
/**
 * 获取图像数据
 * @param url - 图像 URL
 * @param quality - 图像压缩质量（仅当 type 为 'base64' 时有效），范围 0-1，1 表示最高质量
 * @returns - 包含图像数据的 Promise 对象
 */
export async function getImageData(url: string, quality = 1, canvas?: HTMLCanvasElement): Promise<IImageData> {
  const result: IImageData = {
    url,
    blob: null,
    base64: '',
  };
  try {
    const img = await loadImage(url);
    if(!canvas) {
      canvas = document.createElement('canvas');
    }
    canvas.width = img.width;
    canvas.height = img.height;

    const context = canvas.getContext('2d');
    context?.drawImage(img, 0, 0);

    // 获取文件后缀
    const fileExtension = url.split('?')[0].split('.').pop();

    // 获取对应的 MIME 类型
    const mimeType = getMimeTypeFromExtension(fileExtension || '');

    // blob
    result.blob = await new Promise((resolve) => {
      canvas.toBlob((blob) => resolve(blob), mimeType, quality);
    });
    // base64
    result.base64 = canvas.toDataURL(mimeType, quality);
    
    if(context) {
      // 清空画布
      context.fillRect(0, 0, img.width, img.height);
      // 绘制图片
      context.drawImage(img, 0, 0);
      // 清除画布
      context.clearRect(0, 0, img.width, img.height);
    }
  } catch (error) {
    console.error('Error loading or processing image:', error);
  }
  return result;
}
