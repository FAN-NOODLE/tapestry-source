import { IUploadOptions, IUploadResult } from '@/types/upload';
import { cosUploader } from '@/utils/upload';
import { getFileData, getVideoDuration } from '@/utils/upload/utils';
import Api from "@/api/enterprise/material-library";
import { Loading } from 'element-ui';
import { getImageWh } from '@/utils/get-image-info';

// 上传结果扩展接口，增加媒体类型区分
export interface ILocalUploadResult extends IUploadResult {
  /** 媒体类型：image/video */
  mediaType: 'image' | 'video';
}

// 本地上传配置接口
export interface ILocalUploadConfig extends IUploadOptions {
  /** 是否支持多选 */
  multiple?: boolean;
  /** 最大文件大小（字节） */
  maxSize?: number;
  /** 最大文件数量 */
  maxCount?: number;
  /** 允许的文件类型，如 ['image/jpeg', 'video/mp4'] */
  acceptTypes?: string[];
  /** 自定义验证函数，返回false时会拒绝上传 */
  validate?: boolean | ((file: File) => boolean);
  /** 是否进行图片压缩 */
  compress?: boolean;
}

/**
 * 本地上传Mixin
 * 支持图片/视频上传，自动区分文件类型
 */
export default {
  methods: {
    /**
     * 打开文件选择器并上传文件
     * @param config 上传配置
     * @returns Promise 解析为上传结果或上传结果数组（多选时）
     */
    async uploadMedia(config?: ILocalUploadConfig): Promise<ILocalUploadResult | ILocalUploadResult[]> {
      // 创建文件选择input
      const input = document.createElement('input');
      input.type = 'file';
      // 支持多选
      // 最大文件数量
      if (config?.maxCount && config.maxCount > 1) {
        input.multiple = true;
      } else {
        input.multiple = false;
      }
      
      // 设置文件类型过滤
      if (config?.acceptTypes) {
        input.accept = config.acceptTypes.join(',');
      } else {
        // 默认支持常见的图片和视频格式
        input.accept = 'image/jpeg,image/png,video/mp4,video/m4v';
      }
      
      // 设置是否支持多选
      input.multiple = !!config?.multiple;
      
      // 处理文件选择和上传
      return new Promise((resolve, reject) => {
        input.onchange = async (event) => {
          try {
            const files = Array.from((event.target as HTMLInputElement).files || []);
            
            if (files.length === 0) {
              reject(new Error('未选择文件'));
              return;
            }
            
            // 验证文件
            // 收集所有不符合要求的文件
            const oversizedFiles: string[] = [];
            const unsupportedFiles: string[] = [];
            const qyUnsupportedFiles: string[] = [];

            for (const file of files) {
              // 验证文件大小
              if (config?.maxSize && file.size > config.maxSize) {
                oversizedFiles.push(file.name);
              }
              
              // 验证文件类型
              if (config?.acceptTypes && !config.acceptTypes.includes(file.type)) {
                unsupportedFiles.push(file.name);
              }

              // 企业微信不支持JPEG的图片
              if(file.name.toLocaleLowerCase().includes('jpeg')){
                qyUnsupportedFiles.push(file.name);
              }
            }

            // 如果有文件大小超限
            if (oversizedFiles.length > 0) {
              // loading.close();
              reject(new Error(`文件【${oversizedFiles.join('】、【')}】大小超过限制 ${(config.maxSize / 1024 / 1024).toFixed(0)}MB`));
              return;
            }

            // 如果有文件类型不支持
            if (unsupportedFiles.length > 0) {
              reject(new Error(`文件【${unsupportedFiles.join('】、【')}】类型不支持`));
              return;
            }

            // 如果有企业微信不支持的文件
            if (qyUnsupportedFiles.length > 0) {
              reject(new Error(`【${qyUnsupportedFiles.join('】、【')}】文件格式JPEG企业微信不支持`));
              return;
            }
            
            // 验证文件数量
            if (config?.maxCount && files.length > config.maxCount) {
              reject(new Error(`可选文件数量不能超过 ${config.maxCount} 个`));
              return;
            }
            
            // 上传文件
            const uploadPromises = files.map(async (file) => {
              const uploadResult = await cosUploader.uploadFile(file, {
                ...config,
                folder: config?.folder || 'local-upload',
              });
              
              // 判断媒体类型
              let mediaType: 'image' | 'video' = 'image';
              if (file.type.startsWith('video/')) {
                mediaType = 'video';
              }
              // 如果是视频获取封面图
              if (mediaType === 'video') {
                // 获取视频时长
                const duration = await getVideoDuration(file)
                // 视频封面
                const buffer:any = await getFileData(`${uploadResult.url}?ci-process=snapshot&time=1&format=jpg`)
                // 上传视频封面
                const res2:any = await cosUploader.uploadFile(buffer, { filename: uploadResult.url.split('/').pop().split('.')[0] + '_poster.jpg'})
                
                uploadResult.duration = duration
                uploadResult.cover_url = res2.url
              }
              
              return {
                ...uploadResult,
                mediaType,
              } as ILocalUploadResult;
            });
            
            const results = await Promise.all(uploadPromises);
            
            // 如果是单选，返回单个结果；如果是多选，返回结果数组
            resolve(config?.multiple ? results : results[0]);
          } catch (error) {
            reject(error);
          } finally {
            // 清理input元素
            // loading.close();
            input.remove();
          }
        };

        input.oncancel = () => {
          this.loading.close();
        }
        
        // 触发文件选择对话框
        input.click();
      });
    },
    
    /**
     * 上传单个图片
     * @param config 上传配置
     * @returns Promise 解析为图片上传结果
     */
    async uploadImage(config?: Omit<ILocalUploadConfig, 'multiple'>): Promise<ILocalUploadResult> {
      const result = await this.uploadMedia({
        ...config,
        acceptTypes: config?.acceptTypes || ['image/jpeg', 'image/png'],
        multiple: false,
      });

      
      // 确保返回单个结果
      return Array.isArray(result) ? result[0] : result;
    },
    
    /**
     * 上传单个视频
     * @param config 上传配置
     * @returns Promise 解析为视频上传结果
     */
    async uploadVideo(config?: Omit<ILocalUploadConfig, 'multiple'>): Promise<ILocalUploadResult> {
      const result = await this.uploadMedia({
        ...config,
        acceptTypes: config?.acceptTypes || ['video/mp4', 'video/m4v'],
        multiple: false,
      });
      
      // 确保返回单个结果
      return Array.isArray(result) ? result[0] : result;
    },
    
    /**
     * 上传多个媒体文件
     * @param config 上传配置
     * @returns Promise 解析为上传结果数组
     */
    async uploadMedias(config?: ILocalUploadConfig): Promise<ILocalUploadResult[]> {
      this.loading = this.$loading({
        lock: true,
        text: '上传中...',
        background: 'rgba(0, 0, 0, 0.7)',
      });
      try {
        const result = await this.uploadMedia({
          ...config,
          multiple: true,
          // 设置默认最大数量为10个，如果用户没有指定
          maxCount: config?.maxCount || 10,
        });
        
        const list = (Array.isArray(result) ? result : [result])
  
        // 文件压缩
        if(config?.compress){
          await Promise.all(list.map(async (item) => {
            if(item.mediaType === 'image'){
              const { width, height } = await getImageWh(item.url);
              if(width * height > 1080 * 1440){
                const res:any = await Api.compressImage({
                  file_type: 'image',
                  file_path: item.url,
                })
                if (res.code === 200) {
                  item.url = res.data;
                } else {
                  this.$message.error(`素材【${item.name}】压缩失败，请重新选择`);
                }
              }
            }
          }))
        }
  
        const validResults = [];
        if(config?.validate){
          // 验证文件是否符合朋友圈的要求
          await Promise.all(list.map(async (item) => {
            const res:any = await Api.getMediaIdApi({
              url: item.url,
              source_type: item.mediaType === 'image' ? 2 : 5,
              type: 1
            })
    
            if(res.code != 200){
              this.$message.error(`素材【${item.name}】不符合要求，请重新选择`);
              return;
            }
            validResults.push(item);
            
          }))
        } else {
          validResults.push(...list);
        }
  
        // 关闭loa
        
        // 确保返回数组
        return validResults;
        
      } catch (error) {
        console.log('多文件上传失败:', error);
        this.$message.error(error.message);
        throw error;
      } finally {
        // 清理input元素
        this.loading.close();
      }
      
    },
  },
};