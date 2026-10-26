<?php

namespace App\Library\TencentYun\Cos;

use App\Exceptions\Error;
use Laravel\Lumen\Application;
use Qcloud\Cos\Client;

class QCloudCosSdk
{
    /**
     * 腾讯云上传对象实例
     *
     * @var Client
     */
    protected $cosClient;

    /**
     * 存储桶名称
     *
     * @var mixed
     */
    protected $bucket;

    /**
     * 环境域名
     *
     * @var mixed
     */
    protected $host;

    /**
     * 协议头部
     *
     * @var mixed
     */
    protected $schema;

    /**
     * QCloudCos constructor.
     *
     * @param string $bucket
     * @param string $hostUrl
     */
    public function __construct(string $bucket = '', string $hostUrl = '')
    {
        $config          = $this->getConfig();
        $this->cosClient = new Client([
            'region'      => $config['region'],
            'schema'      => $config['schema'],
            'credentials' => [
                'secretId'  => $config['secretId'],
                'secretKey' => $config['secretKey']
            ]
        ]);
        $this->bucket    = empty($bucket) ? $config['bucket'] : $bucket;
        $this->host      = empty($hostUrl) ? $config['host_url'] : $hostUrl;
        $this->schema    = $config['schema'];
    }

    /**
     * 读取配置文件
     *
     * @return array|Application
     */
    private function getConfig()
    {
        $config = config('third_party.QCloudCos');

        if (empty($config) || !is_array($config)) {
            throw new Error(100000, 'QCloudCos配置文件错误，请检查配置文件');
        }

        return $config;
    }

    /**
     * 上传文件
     *
     * @param object $file 文件资源
     * @param string $hash
     * @param string $fileKey
     * @param string $filePath
     *
     * @return string
     */
    public function upload($file, string $hash = 'default', string $fileKey = '', string $filePath = ''): string
    {
        if (empty($fileKey) && empty($filePath)) {
            $filePath  = $file->path();
            $md5Uniqid = md5(uniqid());
            $picName   = $hash . '/' . date('Y/m/d/', time()) . $md5Uniqid;
            $fileKey   = $picName . '.' . $file->getClientOriginalExtension();
        }

        try {
            $result = $this->cosClient->putObject([
                'Bucket' => $this->bucket,
                'Key'    => $fileKey,
                'Body'   => fopen($filePath, 'rb'),
            ])->toArray();
            $url    = $this->schema . "://" . $result['Location'];
        } catch (\Exception $exception) {
            throw new Error(100000, '图片上传失败');
        }

        return $url;
    }

    /**
     * 通过文件内容
     *
     * @param $content
     *
     * @return string
     */
    public function uploadContent($content): string
    {

        $key = date('Y/m/d/', time()) . md5(uniqid()) . '.png';

        try {
            $result = $this->cosClient->putObject([
                'Bucket' => $this->bucket,
                'Key'    => $key,
                'Body'   => $content,
            ])->toArray();
            $url    = $this->schema . "://" . $result['Location'];
        } catch (\Exception $exception) {
            throw new Error(100000, '文件上传失败');
        }

        return $url;
    }


    /**
     * 管理后台下载中心文件上传
     *
     * @param string $filePath 文件地址
     *
     * @return mixed
     */
    public function downloadUploadFile($filePath)
    {
        $path = pathinfo($filePath);
        $key  = 'download/' . date('Y/m/d/', time()) . $path['filename'] . '.' . mt_rand(100, 999) . '.' . $path['extension'];

        try {
            $result = $this->cosClient->putObject([
                'Bucket' => $this->bucket,
                'Key'    => $key,
                'Body'   => fopen($filePath, 'rb'),
            ])->toArray();
            $url    = $this->schema . "://" . $result['Location'];
        } catch (\Exception $exception) {
            throw new Error(100000, '文件上传失败');
        }

        return $url;
    }

    /**
     * 通过文件url地址上传文件
     *
     * @param $imageUrl
     *
     * @return string
     */
    public function uploadImageUrl($imageUrl)
    {

        $key = date('Y/m/d/', time()) . md5(uniqid()) . '.png';

        try {
            $result = $this->cosClient->putObject([
                'Bucket' => $this->bucket,
                'Key'    => $key,
                'Body'   => file_get_contents($imageUrl),
            ])->toArray();

            $url = $this->schema . "://" . $result['Location'];

        } catch (\Exception $exception) {
            throw new Error(100000, '图片上传失败');
        }

        return $url;
    }

    /**
     * 返回腾讯云签名
     *
     * @return string
     */
    public function getSign()
    {

        // 确定 App 的云 API 密钥
        $config     = $this->getConfig();
        $secret_id  = $config['secretId'];
        $secret_key = $config['secretKey'];

        // 确定签名的当前时间和失效时间
        $current = time();
        $expired = $current + 86400;  // 签名有效期：1天

        // 向参数列表填入参数
        $arg_list = array(
            "secretId"         => $secret_id,
            "currentTimeStamp" => $current,
            "expireTime"       => $expired,
            "random"           => rand());

        // 计算签名
        $orignal   = http_build_query($arg_list);
        $signature = base64_encode(hash_hmac('SHA1', $orignal, $secret_key, true) . $orignal);

        return $signature;
    }
}
