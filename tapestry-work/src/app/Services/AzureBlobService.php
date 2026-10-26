<?php
/**
 * azure blob 服务层
 *
 * User: wangle@vchangyi.com
 * Date: 2021/03/08
 */

namespace App\Services;


use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\Storage;

class AzureBlobService
{
    // 配置信息
    protected $blobConfig = [];

    // 存储路径
    protected $saveBlobDir = '';


    /**
     * AzureBlobService constructor.
     *
     */
    public function __construct()
    {
        $config = config('filesystems.disks.azure');
        $this->blobConfig = [
            'azureBlobUrl'          => 'https://' . $config['name'] . '.blob.core.chinacloudapi.cn',
            'azureStorageContainer' => $config['container'],
            'azureStorageAccount'   => $config['name'],
            'azureStorageKey'       => $config['key'],
        ];

        $this->saveBlobDir = env("APP_ENV").'_'.env("APP_NAME").'_'.date('Ymd');

    }


    /**
     * 文件上传
     * @param $file
     * @return string|null
     */
    public function uploadByFile($file)
    {
        if(!empty($file)){
            $path = Storage::disk('azure')->put($this->saveBlobDir,$file);
            $url  = $this->blobConfig['azureBlobUrl'] . '/' . $this->blobConfig['azureStorageContainer'] . '/' . $path;
            return $url;
        }
        return null;

    }


    /**
     * url上传
     * @param $url
     * @return string|null
     */
    public function uploadByUrl($url)
    {
        $urlInfo      = pathinfo($url);
        $tempFileName = time() . md5(uniqid()) . '.' . (isset($urlInfo['extension']) ? $urlInfo['extension'] : 'png');

        $bool         = Storage::disk('local')->put($tempFileName,file_get_contents($url));
        if($bool){
            $tempPath = config('filesystems.disks.local.root') . '/' . $tempFileName;
            $file     = new UploadedFile($tempPath,$urlInfo['basename']);
            $path     = Storage::disk('azure')->put($this->saveBlobDir,$file);
            if($path){
                @unlink($tempPath);
                $url  = $this->blobConfig['azureBlobUrl'] . '/' . $this->blobConfig['azureStorageContainer'] . '/' . $path;
                return $url;
            }

        }
        return null;

    }


    /**
     * 文件上传
     * @param $filePath
     * @return string|null
     */
    public function uploadByFilePath($filePath)
    {
        if(is_file($filePath)){
            $filePathInfo = pathinfo($filePath);
            $file         = new UploadedFile($filePath,$filePathInfo['basename']);

            $path = Storage::disk('azure')->put($this->saveBlobDir,$file);
            if($path){
                $url = $this->blobConfig['azureBlobUrl'] . '/' . $this->blobConfig['azureStorageContainer'] . '/' . $path;
                return $url;
            }

        }

        return null;
    }


    /**
     * @param $content
     * @return string|null
     */
    public function uploadByContent($content)
    {
        $tempFileName = time() . md5(uniqid()) . '.png';
        $bool         = Storage::disk('local')->put($tempFileName, $content);
        if ($bool) {
            $tmpPath  = config('filesystems.disks.local.root').'/'.$tempFileName;
            $file     = new UploadedFile($tmpPath, $tempFileName);
            $path     = Storage::disk('azure')->put($this->saveBlobDir, $file);
            if ($path) {
                @unlink($tmpPath);
                $url = $this->blobConfig['azureBlobUrl'] . '/' . $this->blobConfig['azureStorageContainer'] . '/' . $path;
                return $url;
            }
        }
        return null;
    }


    /**
     * 管理后台下载中心文件上传
     *
     * @param string $filePath 文件地址
     * @return mixed
     */
    public function downloadUploadByFile($filePath)
    {
        $urlInfo      = pathinfo($filePath);
        $tempFileName = 'download/' . time() . md5(uniqid()) . '.' . (isset($urlInfo['extension']) ? $urlInfo['extension'] : 'png');

        $bool         = Storage::disk('local')->put($tempFileName, $filePath);
        if($bool){
            $tempPath = config('filesystems.disks.local.root') . '/' . $tempFileName;
            $file     = new UploadedFile($tempPath,$urlInfo['basename']);
            $path     = Storage::disk('azure')->put($this->saveBlobDir,$file);
            if($path){
                @unlink($tempPath);
                $url  = $this->blobConfig['azureBlobUrl'] . '/' . $this->blobConfig['azureStorageContainer'] . '/' . $path;
                return $url;
            }

        }
        return null;

    }

}