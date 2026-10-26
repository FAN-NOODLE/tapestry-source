/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import html2canvas from '@whowasme/html2canvas';
import htmlToPdfmake from 'html-to-pdfmake';
import pdfmake from "pdfmake/build/pdfmake";
import pdfFonts from "pdfmake/build/vfs_fonts";
import { getFileBase64 } from '.';
import { getImageData } from './util';
import { cosUploader } from '../upload';
pdfmake.vfs = pdfFonts.pdfMake.vfs;

function loadFonts(list: string[]) {
  // const url = window.location.origin + '/work';
  return Promise.all(list.map((item) => getFileBase64(`${item}`))).then((res) => {
    list.forEach((item, index) => {
      const key = item.split('/').pop() as string;
      pdfmake.vfs[key] = res[index] as string;
    });
  });
}

const list = [
  'https://cdn-cos.coach.com.cn/uploads/fonts/STSongti-SC-Bold.ttf',
  'https://cdn-cos.coach.com.cn/uploads/fonts/STSongti-SC-Regular.ttf',
  'https://cdn-cos.coach.com.cn/uploads/fonts/Arial.ttf',
  'https://cdn-cos.coach.com.cn/uploads/fonts/ArialBold.ttf',
  'https://cdn-cos.coach.com.cn/uploads/fonts/ArialItalic.ttf',
  'https://cdn-cos.coach.com.cn/uploads/fonts/ArialBoldItalic.ttf',
];

export interface ISavePdfFromHtmlOptions {
  element: HTMLElement;
  filename?: string;
  type?: string;
  background?: (currentPage:number, pageSize:number) => any;
  footer?: (currentPage:number, pageSize:number) => any;
  watermark?: any,
  styles?: {
    [key: string]: any;
  };
  info?: {
    title?: string;
    author?: string;
    subject?: string;
    keywords?: string;
    creator?: string;
    producer?: string;
    creationDate?: Date;
    modDate?: Date;
  };
}

interface IProgressCallback {
  (percent: {
    percent: number;
    status: 'loading'|'saving'|'success'|'error';
  }): void;
}

export class SavePdfFromHtml {
  options: ISavePdfFromHtmlOptions;
  allImageMap: Map<string, string>;
  element: HTMLElement;
  onProgress: IProgressCallback;
  constructor(options: ISavePdfFromHtmlOptions) {
    this.options = {
      filename: 'pdf.pdf', 
      type: 'pdf', 
      ...options
    };
    this.element = options.element;
    this.allImageMap = new Map();
    this.onProgress = () => {
      // do nothing
    };
  }

  async create() {
    this.onProgress({
      percent: 0,
      status: 'loading'
    });
    // 加载字体
    await loadFonts(list)
    // 配置字体
    pdfmake.fonts = {
      Arial: {
        normal: 'Arial.ttf',
        bold: 'ArialBold.ttf',
        italics: 'ArialItalic.ttf',
        bolditalics: 'ArialBoldItalic.ttf',
      },
      SongTiSC: {
        normal: 'STSongti-SC-Regular.ttf',
        bold: 'STSongti-SC-Bold.ttf',
        italics: 'STSongti-SC-Regular.ttf',
        bolditalics: 'STSongti-SC-Bold.ttf',
      },
    };

    // 获取所有资源
    this.getAssets()
    // 加载资源
    await this.loadAssets();
    // 替换资源
    await this.replaceAssets();
    // 表格转图片
    await this.tableToImage()
    
    this.onProgress({
      percent: 0.6,
      status: 'loading'
    });
    const str = this.replaceChinese(this.element.innerHTML)
    this.onProgress({
      percent: 0.8,
      status: 'loading'
    });
    // console.log(str);
    
    const content = htmlToPdfmake(str, {
      tableAutoSize: true
    })
    this.onProgress({
      percent: 0.9,
      status: 'loading'
    });
    const { filename, styles, info, background, footer, watermark } = this.options

    const params:any = {
      content,
      styles,
      defaultStyle: {
        font: 'Arial',
        lineHeight: 1,
      },
      info: info,
      footer,
      background,
      watermark
    }
    console.log(content);
    // debugger
    

    // pdfmake.createPdf(params).download(filename);
    return new Promise((resolve, reject) => {
      pdfmake.createPdf(params).getBlob((res) => {
        cosUploader.uploadFile(res, { filename, folder:'uploads/pdf' }).then((res: any) => {
          this.onProgress({
            percent: 1,
            status: 'success'
          });
          resolve({
            name: filename,
            url: res.url
          });
        }).catch(err => {
          reject(err)
        })
      });
    })
    
  }

  // 匹配任何中文字符  
  replaceChinese(str: string) {  
    const symbolRegex = /[，。！？、；：“”‘’（）《》【】〈〉（）《》「」『』﹝﹞｛｝（）《》〔〕【】〈〉《》「」『』︵︶︷︸︹︺︿﹀﹁﹂﹃﹄﹈﹉﹊﹋﹌﹍﹎﹏﹐﹑﹒﹔﹕﹖﹗﹙﹚﹛﹜﹝﹞（）《》〔〕｛｝〈〉｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝﾞﾟ]/g;  
    const chineseRegex = /[\u4e00-\u9fa5]/g;  
    const regex = new RegExp(`${symbolRegex.source}|${chineseRegex.source}`, 'g')
    const res =  str.replace(regex, ($1) => {
      return `<span class="zh-font">${$1}</span>`;
    });  

    // 匹配a标签
    const aRegex = /<a(.*?)>(.*?)<\/a>/g;
    const aRes = res.replace(aRegex, ($1, $2, $3) => {
      const ctx1 = $1.replaceAll('<span class="zh-font">', '').replaceAll('</span>', '');
      return ctx1;
    });

    // 替换所有的padding
    const paddingRegex = /padding-left:(.*?);/g;
    const paddingRes = aRes.replace(paddingRegex, ($1, $2) => {
      return `margin-left:${$2};`;
    });
    return paddingRes
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
      const percent = parseFloat((((i+1)/(len))).toFixed(2)) * 0.4;
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
      const width = Math.min(image.width, 600);
      const height = image.height * (width / image.width);
      
      image.setAttribute('width', width + 'px')
      // image.setAttribute('height', height + 'px')
      
      // 获取图片元素的 src 属性
      const src = image.getAttribute('src') || '';
  
      // 对 src 进行替换，得到新的图片地址
      // 如果图片是远程地址，则使用 imageToBase64 方法将其转换为 base64
      const isUrl = /^http(s?):\/\//i.test(src);
      if(isUrl) {
        this.allImageMap.set(src, '');
      }
    }

    // 获取所有a标签
    const aList = document.getElementsByTagName('a');
    for (let i = 0; i < aList.length; i++) {
      const a = aList[i];
      a.className = ['zh-font', a.className].filter(Boolean).join(' ')
    }
  }


  async tableToImage() {
    // 获取main下面的所有表格
    const mainElements = this.element.querySelectorAll('.main table');
    
    // const tableElements = []
    for (let i = 0; i < mainElements.length; i++) {
      const table = mainElements[i] as HTMLTableElement;
      table.style.width = table.clientWidth + 'px'
      // 删除属性
      table.removeAttribute('width')
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

  setOptions(options: ISavePdfFromHtmlOptions) {
    this.options = {
      ...this.options,
      ...options
    }
    return this
  }

  progress(callback: IProgressCallback) {
    this.onProgress = callback.bind(this)
    return this
  }
}