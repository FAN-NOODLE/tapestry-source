
import html2canvas, { Options } from '@whowasme/html2canvas';
import { download } from '.';
import { getImageData } from './util';

export interface ISaveImageFromHtmlOptions extends Partial<Options> {
  element: HTMLElement;
  filename?: string;
  type?: string;
  quality?: number;
} 

interface IProgressCallback {
  (percent: {
    percent: number;
    status: 'loading'|'saving'|'success'|'error';
  }): void;
}

export class SaveImageFromHtml {
  options: ISaveImageFromHtmlOptions;
  canvas: HTMLCanvasElement;
  allImageMap: Map<string, string>;
  element: HTMLElement;
  onProgress: IProgressCallback;
  constructor(options: ISaveImageFromHtmlOptions) {
    this.options = {
      filename: 'image.png', 
      type: 'image/png', 
      quality: 1,
      useCORS: true,
      ...options
    };
    this.element = options.element;
    this.canvas = document.createElement('canvas');
    this.allImageMap = new Map();
    this.onProgress = () => {
      // do nothing
    };
  }

  async create(autosave = true) {
    const { filename = 'image', type = 'image/png', quality = 1, ...options } = this.options;
    this.onProgress({
      percent: 0,
      status: 'loading'
    });
    // 获取所有资源
    this.getAssets()
    // 加载资源
    await this.loadAssets();
    // 替换资源
    await this.replaceAssets();
    
    this.onProgress({
      percent: 0.9,
      status: 'saving'
    });
    // 绘制canvas
    const canvas = await html2canvas(this.element, options);
  
    // base64
    const base64 = canvas.toDataURL(type, quality);
    // blob
    const blob: Blob = await new Promise((resolve) => {
      canvas.toBlob((blob) => resolve(blob), type, quality);
    });
    
    // 自动下载
    if(autosave) {
      download(blob, type, filename)
    }

    this.onProgress({
      percent: 1,
      status: 'success'
    });

    return {
      base64,
      blob
    }
  }

  async loadAssets() {
    const assets = [...this.allImageMap.keys()]
    const len = assets.length;
    for (let i = 0; i < len; i++) {
      const url = assets[i];
      // 获取图片数据
      const { base64 } = await getImageData(url)
      // 记录数据
      this.allImageMap.set(url, base64) 
      // 显示进度
      const percent = parseInt((((i+1)/(len))).toFixed(2)) * 0.9;
      this.onProgress({
        percent,
        status: 'loading'
      });
    }
  }

  /**
   * @description 获取所有图片
   */
  getAssets() {
    
    // 获取页面上所有的图片元素
    const images = this.element.getElementsByTagName('img');
    
    for (let i = 0; i < images.length; i++) {
      const image = images[i];
      // 获取图片元素的 src 属性
      const src = image.getAttribute('src') || '';
  
      // 对 src 进行替换，得到新的图片地址
      // 如果图片是远程地址，则使用 imageToBase64 方法将其转换为 base64
      const isUrl = /^http(s?):\/\//i.test(src);
      if(isUrl) {
        this.allImageMap.set(src, '');
      }
    }
  }

  /**
   * 替换所有元素的src
   * @param element 
   */
  replaceAssets() {
    // 获取页面上所有的图片元素
    return new Promise((resolve) => {
      const images = this.element.getElementsByTagName('img');
      for (let i = 0; i < images.length; i++) {
        const image = images[i];
        // 获取图片元素的 src 属性
        const src = image.getAttribute('src') || '';

        // 将图片元素的 src 属性替换为新地址
        const base64 = this.allImageMap.get(src)
        if(base64) {
          image.setAttribute('src', base64);
        }
      }
      resolve(1)
    });
  }

  setOptions(options: ISaveImageFromHtmlOptions) {
    this.options = {
      ...this.options,
      ...options
    }
    return this
  }

  progress(callback: IProgressCallback) {
    this.onProgress = callback
    return this
  }
}