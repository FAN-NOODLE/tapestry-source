import COS from 'cos-js-sdk-v5';
import { getFilePath } from './utils';
import { getFileCredentialApi } from '@/api/marketing-promotion/index';
const uat = 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/'
const uatcdn = 'https://cdn-cos-uat.coach.com.cn/'
const prod = 'https://tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com/'
const prodcdn = 'https://cdn-cos.coach.com.cn/'
const publicPrefixes = ['local-upload/', 'static/', 'uploads/']


export default class CosUploader {
  private client: COS;
  private credentialReady: Promise<void>;
  credential: any;

  constructor(private readonly config: COS.COSOptions) {
    this.client = new COS(config);
    this.credentialReady = this.init()
  }
  async init() {
    const res:any = await getFileCredentialApi();
    this.credential = res.data;
  }

  private normalizeKey(filePath: string): string {
    return (filePath || '').replace(/^\/+/, '');
  }

  private isPublicKey(filePath: string): boolean {
    const key = this.normalizeKey(filePath);
    return publicPrefixes.some(prefix => key.startsWith(prefix));
  }

  private getBucket(filePath: string): string {
    return this.isPublicKey(filePath) ? this.credential?.bucket : (this.credential?.private_bucket || this.credential?.bucket);
  }

  private formatFileUrl(location: string, filePath: string): string {
    const url = `https://${location}`;
    if (this.isPublicKey(filePath)) {
      return url.replace(uat, uatcdn).replace(prod, prodcdn);
    }
    return url;
  }

  private getObjectInfo(fileUrl: string): { bucket: string; key: string } {
    const decodedUrl = decodeURIComponent(fileUrl);
    const privateBucket = this.credential?.private_bucket;
    const publicOrigin = `https://${this.credential?.bucket}.cos.${this.credential?.region}.myqcloud.com/`;
    const privateOrigin = privateBucket ? `https://${privateBucket}.cos.${this.credential?.region}.myqcloud.com/` : '';
    if (privateOrigin && decodedUrl.startsWith(privateOrigin)) {
      return { bucket: privateBucket, key: decodedUrl.replace(privateOrigin, '') };
    }
    return {
      bucket: this.credential?.bucket,
      key: decodedUrl.replace(publicOrigin, '').replace(uatcdn, '').replace(prodcdn, ''),
    };
  }

  uploadFile(file: any, config?: IUploadOptions) {
    return new Promise((resolve, reject) => {
      const isFile = file instanceof File
      console.log(isFile);
      
      const filePath = isFile ? getFilePath(file, config?.folder) :getFilePath(config?.filename, config?.folder);
      const bucket = this.getBucket(filePath);
      const formatFileUrl = (location: string) => this.formatFileUrl(location, filePath);
      this.client.uploadFile(
        {
          Bucket: bucket,
          Region: this.credential?.region,
          Key: filePath,
          Body: file, // 上传文件对象
          SliceSize: 1024 * 1024 * 2 /* 触发分块上传的阈值，超过10MB使用分块上传，非必须 */,
          Headers: {
            'Content-Disposition': `attachment; filename="${encodeURIComponent(config?.filename||'')}"`
          },
          onProgress: (progressData: { percent: number }) => {
            config?.onProgress?.(progressData); //上传进度
          },
        },
        function (err, data) {
          if (err) {
            reject(err);
          } else {
            if(isFile) {  
              resolve({
                name: file.name,
                size: file.size,
                type: file.type,
                ext: file.name.split('.').pop() || '',
                url: formatFileUrl(data.Location),
              })
            } else {
              resolve({
                url: formatFileUrl(data.Location),
              });
            }
          }
        },
      )
    })
  }
  sliceUploadFile(file: any, config?: IUploadOptions) {
    return new Promise((resolve, reject) => {
      const isFile = file instanceof File
      const filePath = isFile ? getFilePath(file, config?.folder) :getFilePath(config?.filename, config?.folder);
      const bucket = this.getBucket(filePath);
      const formatFileUrl = (location: string) => this.formatFileUrl(location, filePath);
      this.client.sliceUploadFile(
        {
          Bucket: bucket,
          Region: this.credential?.region,
          Key: filePath,
          Body: file, // 上传文件对象
          // SliceSize: 1024 * 1024 * 10 /* 触发分块上传的阈值，超过10MB使用分块上传，非必须 */,
          Headers: {
            'Content-Disposition': `attachment; filename="${encodeURIComponent(config?.filename||'')}"`
          },
          onProgress: (progressData: { percent: number }) => {
            config?.onProgress?.(progressData); //上传进度
          },
        },
        function (err, data) {
          if (err) {
            reject(err);
          } else {
            if(isFile) {  
              resolve({
                name: file.name,
                size: file.size,
                type: file.type,
                ext: file.name.split('.').pop() || '',
                url: formatFileUrl(data.Location),
              })
            } else {
              resolve({
                url: formatFileUrl(data.Location),
              });
            }
          }
        },
      )
    })
  }
  async getObjectUrl(url: string): Promise<string> {
    await this.credentialReady;
    // 去除链接域名
    const objectInfo = this.getObjectInfo(url);
    return new Promise((resolve, reject) => {
      this.client.getObjectUrl(
        {
          Bucket: objectInfo.bucket /* 必须 */,
          Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
          Key: objectInfo.key, // 存储在桶里的对象键（例如1.jpg，a/b/test.txt），支持中文，必须字段
        },
        function (err, data) {
          console.log(err || data.Url);
          if(err) {
            reject(err);
          } else {
            resolve(data.Url);
          }
        }
      );
    })
  }
}

export const cosUploader = new CosUploader({
  getAuthorization: async function (options, callback) {
    const res:any = await getFileCredentialApi();
    if (res.code !== 200) throw new Error('获取临时密钥失败');
    callback({
      TmpSecretId: res.data?.secret_id,
      TmpSecretKey: res.data?.secret_key,
      XCosSecurityToken: res.data?.session_token,
      StartTime: res.data?.start_time,
      ExpiredTime: res.data?.expired_time,
    });
  },
});
