import JSZip from 'jszip';
import { download, downloadFile, getFileBlob } from '@/utils/download';
import { getImageData } from './util';
import { cosUploader } from '../upload';


export interface IToZipItem {
  name: string;
  content: string | Blob;
}

/**
 * 保存图片到 zip
 * @param data 
 * @param filename 
 */
export async function saveImageToZip(data: IToZipItem[], filename='example.zip') {
  const zip = new JSZip();

  const promise = data.map(async (item) => {
    if(typeof item.content ==='string' && item.content.startsWith('http')) {
      return {
        name: item.name,
        content: await getImageData(item.content)
      }
    } else {
      return item
    }
  })

  const result = await Promise.all(promise);

  result.forEach(item => {
    zip.file(item.name, item.content)
  })


  // 生成 Zip 文件
  const blob = await zip.generateAsync({ type: "blob" })
  download(blob, '', filename);

}
/**
 * 保存文件到 zip
 * @param data 
 * @param filename 
 */
export async function saveFileToZip(data: IToZipItem[], filename='example.zip') {
  const zip = new JSZip();

  const promise = data.map(async (item) => {
    if(typeof item.content ==='string' && item.content.startsWith('http')) {
      const url = await cosUploader.getObjectUrl(item.content)
      // 判断链接是图片或视频
      const fileData = (await getFileBlob(url)).data
      return {
        name: item.name,
        content: fileData
      }
    } else {
      return item
    }
  })

  const result = await Promise.all(promise);


  // 整理数据，处理相同名称的文件
  const fileMap = new Map<string, IToZipItem[]>()

  function setFileMap(item, index=1) {
    // 重命名
    const names = item.name.split('.')
    names[0] += ` (${index})`
    const name = names.join('.')
    
    if(fileMap.has(name)) {
      setFileMap(item, index+1)
    } else {
      item.name = name
      fileMap.set(item.name, [item])
    }
  }

  result.forEach(item => {
    if(fileMap.has(item.name)) {
      setFileMap(item)
    } else {
      fileMap.set(item.name, [item])
    }
  })
  
  fileMap.forEach((value, key) => {
    zip.file(value[0].name, value[0].content)
  })


  // 生成 Zip 文件
  const blob = await zip.generateAsync({ type: "blob" })
  download(blob, '', filename);

}

