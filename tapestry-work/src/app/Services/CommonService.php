<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2020/2/3                   
* +------------------------------------+
* | Time   : 13:33                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Library\Uc\EnterprisePlugin;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use CrsSdk\Boss\ReleaseEnterpriseWork;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;
use Qcloud\Cos\Client;

class CommonService
{
    private static $publicCosPrefixes = [
        'local-upload/',
        'static/',
        'uploads/',
    ];

    private static function normalizeCosKey($key)
    {
        return ltrim($key, '/');
    }

    private static function isPublicCosKey($key)
    {
        $key = self::normalizeCosKey($key);
        foreach (self::$publicCosPrefixes as $prefix) {
            if (strpos($key, $prefix) === 0) {
                return true;
            }
        }
        return false;
    }

    private static function getCosBucket($key, $config)
    {
        if (self::isPublicCosKey($key)) {
            return $config['bucket'];
        }
        return $config['private_bucket'] ?? $config['bucket'];
    }

    private static function getCosUrl($key, $config)
    {
        $key = self::normalizeCosKey($key);
        if (self::isPublicCosKey($key)) {
            return $config['schema'] . '://' . rtrim($config['cdn_host_url'], '/') . '/' . $key;
        }
        return $config['schema'] . '://' . rtrim($config['private_host_url'], '/') . '/' . $key;
    }

    private static function defaultCosKey()
    {
        return 'uploads/' . date('Y/m/d/', time()) . md5(uniqid()) . '.png';
    }

    public static function uploadImage($path, $key = '')
    {
        $config = config('cloud.QCloudCos');

        if (empty($config)) {
            throw new Error(90000, 'QCloudCos配置错误');
        }

        $client = new Client([
            'region'      => $config['region'],
            'schema'      => $config['schema'],
            'credentials' => [
                'secretId'  => $config['secretId'],
                'secretKey' => $config['secretKey']
            ]
        ]);

        $key = $key ?: self::defaultCosKey();
        $key = self::normalizeCosKey($key);

        try {
            $client->putObject([
                'Bucket' => self::getCosBucket($key, $config),
                'Key'    => $key,
                'Body'   => fopen($path, 'rb'),
            ]);

        } catch (\Exception $exception) {
            throw new Error(50000, '图片上传异常');
        }

        return self::getCosUrl($key, $config);
    }

    public static function uploadContent($content, $key = '')
    {
        $tmpPath = tempnam(sys_get_temp_dir(), 'cos-content-');
        file_put_contents($tmpPath, $content);

        try {
            return self::uploadImage($tmpPath, $key);
        } finally {
            if (file_exists($tmpPath)) {
                unlink($tmpPath);
            }
        }
    }

    public static function uploadImageUrl($url, $key = '')
    {
        $content = file_get_contents($url);
        if ($content === false) {
            throw new Error(50000, '图片下载异常');
        }

        return self::uploadContent($content, $key);
    }

    /**
     * @param array  $list
     * @param int    $id
     * @param string $pk
     * @param string $pik
     * @param string $field
     *
     * @return array
     */
    public static function getSubs(Array $list = [], $id = 0, $pk = 'dpId', $pik = 'dpParentid', $field = 'dpThirdid')
    {
        $subs = [];
        foreach ($list as $key => $val) {
            if ($val[$pik] == $id) {
                $subs[] = $val[$field];
                unset($list[$key]);
                $subs = array_merge($subs, self::getSubs($list, $val[$pk], $pk, $pik, $field));
            }
        }
        return $subs;
    }

    /**
     * @param $enterpriseHash
     * @return array|mixed
     */
    public function getAuthDetail($enterpriseHash)
    {
        $key = redisKey('REDIS_SHOP_AUTH_PREFIX', ['hash' => $enterpriseHash]);
        $authData = getOrSetData($key);

        if (!$authData) {
            $authData = (new Auth())->authDetail([
                'enterprise_hash' => $enterpriseHash,
                'auth_type'       => Constant::APP_TYPE
            ])['data'];
            getOrSetData($key, $authData);
        }
        return $authData;
    }

    /**
     * 获取门店列表
     * @param $params
     */
    public function getStoreList($params)
    {
        // 防止键名过长
        $key = md5(json_encode($params));
        $key = redisKey('STORE_LIST', ['key' => $key]);

        $list = getOrSetData($key);
        if (!$list) {
            $list = (new Store())->listByMobile($params);

            getOrSetData($list, $list, 600);
        }

        return $list;
    }
}
