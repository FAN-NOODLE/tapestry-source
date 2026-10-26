<?php

/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Modules\ThirdPart;

use App\Exceptions\Error;
use App\Library\Util\ImageUtil;
use App\Models\EnterpriseThirdPartAccount;
use App\Modules\BaseModule;
use App\Providers\ThirdPartFacade;
use CrsSdk\Common\QCloudCos;
use EasyWeChat\Kernel\Exceptions\DecryptException;
use Illuminate\Support\Facades\Cache;

class MiniProgramModule extends BaseModule
{
    /**
     * 获取session key.
     *
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param string $code
     *
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function getSessionKey(EnterpriseThirdPartAccount $thirdPartAccount, string $code)
    {
        /*
         * [
         *  "openid" => "123",
         *  "session_key" => "1132"
         * ]
         */
        // 判断 code 是否缓存对应的 session_key
        /*
        $cacheKey = $thirdPartAccount->type . ':' . $code;
        if (Cache::has($cacheKey)) {
            $res = Cache::get($thirdPartAccount->type . ':' . $code);
        } else {
            $res = ThirdPartFacade::miniProgram($thirdPartAccount->store_id, $thirdPartAccount->type)->auth->session($code);
            if (key_exists('errcode', $res)) {
                throw new Error(100000,'error_msg: ' . $res['errmsg']);
            }
            Cache::put($cacheKey, $res, 12 * 60);
        }
        */

        $res = ThirdPartFacade::miniProgram($thirdPartAccount->enterprise_id, $thirdPartAccount->type)->auth->session($code);
        if (key_exists('errcode', $res)) {
            throw new Error(100000, 'error_msg: ' . $res['errmsg']);
        }

        return $res;
    }

    /**
     * 解密用户信息.
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param $sessionKey
     * @param $iv
     * @param $encryptedData
     * @return array
     * @throws Error
     */
    public function getUserEncryptedData(EnterpriseThirdPartAccount $thirdPartAccount, $sessionKey, $iv, $encryptedData)
    {
        /*
         * [
         *      "openId" => "OPENID",
         *      "nickName" => "NICKNAME",
         *      "gender" => "GENDER",
         *      "city" => "CITY",
         *      "province" => "PROVINCE",
         *      "country" => "COUNTRY",
         *      "avatarUrl" => "AVATARURL",
         *      "unionId" => "UNIONID",
         *      "watermark" =>
         *      [
         *          "appid" => "APPID",
         *          "timestamp" => TIMESTAMP
         *      ]
         * ]
         */
        try {
            $res = ThirdPartFacade::miniProgram($thirdPartAccount->enterprise_id, $thirdPartAccount->type);
            $res->encryptor->decryptData($sessionKey, $iv, $encryptedData);
        } catch (DecryptException $e) {
            throw new Error(100000, '微信code失效，请重新获取');
        }

        if (key_exists('errcode', $res)) {
            throw new Error(100000, 'error_msg: ' . $res['errmsg']);
        }

        return $res;
    }

    /**
     * 小程序主二维码
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param $path
     * @param int $type 1=方形码，2=圆形码
     * @return mixed
     */
    public static function mainMiniProgramQrCode(EnterpriseThirdPartAccount $thirdPartAccount, $path, $type = 1)
    {
        if ($thirdPartAccount->main_qrcode) { // 不存在主二维码就直接生成一个
            return $thirdPartAccount->main_qrcode;
        }
        $miniProgram = ThirdPartFacade::miniProgram($thirdPartAccount->enterprise_id, $thirdPartAccount->type);
        if ($type == 1) {
            $res = $miniProgram->app_code->getQrCode($path);
        } else {
            $res = $miniProgram->app_code->get($path);
        }

        if (in_array('image/jpeg', $res->getHeader('Content-Type'))) {

            $dirName = env('APP_ENV') . '-logs/qrcode';
            $filename = $res->save(storage_path($dirName), $thirdPartAccount->app_id);

            $path = storage_path($dirName) . '/' . $filename;
            $url = (new QCloudCos())->uploadContent(file_get_contents($path));

            $thirdPartAccount->update([
                'main_qrcode' => $url
            ]);
        } else {
            throw new Error('100000', '获取小程序二维码失败,微信返回:' . json_encode($res));
        }

        // 返回路径
        $thirdPartAccount->refresh();
        return $thirdPartAccount->main_qrcode;
    }

    /**
     * 生成自定义小程序二维码
     *
     * @param $path
     * @param $type
     * @param $storeId
     * @param $qr_type 1=方形码，2=圆形码
     * @return null
     * @throws Error
     * @throws \Exception
     */
    public function customQrCode($path, $type, $enterpriseId, $qr_type = 1)
    {
        // 初始化小程序实例
        $app = ThirdPartFacade::miniProgram($enterpriseId, $type);
        if ($qr_type == 1) {
            $res = $app->app_code->getQrCode($path);
        } elseif ($qr_type == 2) {
            $res = $app->app_code->get($path);
        }

        $qrcodeUrl = null;
        $resType = $res->getHeader('Content-Type'); // 用于判断是否是成功返回的二维码图片
        if (in_array('image/jpeg', $resType)) {
            $image = uploadImage($res, ImageUtil::TYPE_QR_CODE);
            $qrcodeUrl = $image->url;
        } else {
            throw new Error('100000', '获取小程序二维码失败，微信返回：' . json_encode($res));
        }

        // 返回路径
        return $qrcodeUrl;
    }

    /**
     * 小程序带场景参数二维码
     *
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param $path
     * @param $query
     * @return mixed
     * @throws Error
     * @throws \Exception
     */
    public function pathMiniProgramQrCode(EnterpriseThirdPartAccount $thirdPartAccount, $path, $query)
    {
        $miniProgram = ThirdPartFacade::miniProgram($thirdPartAccount->enterprise_id, $thirdPartAccount->type);
        $scene = http_build_query($query);

        $res = $miniProgram->app_code->getUnlimit($scene, ['page' => $path]);

        $resType = $res->getHeader('Content-Type'); // 用于判断是否是成功返回的二维码图片
        if (in_array('image/jpeg', $resType)) {
            $image = uploadImage($res, ImageUtil::TYPE_QR_CODE);
        } else {
            throw new Error('100000', '获取小程序二维码失败，微信返回：' . json_encode($res));
        }

        // 返回路径
        return $image->url;
    }

    /**
     * 小程序带场景参数二维码
     *
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param $path
     * @param $scene
     * @return mixed
     * @throws Error
     * @throws \Exception
     */
    public function pathMiniProgramSceneQrCode(EnterpriseThirdPartAccount $thirdPartAccount, $path, $scene = '')
    {
        $miniProgram = ThirdPartFacade::miniProgram($thirdPartAccount->enterprise_id, $thirdPartAccount->type);

        $res = $miniProgram->app_code->getUnlimit($scene, ['page' => $path]);

        $resType = $res->getHeader('Content-Type'); // 用于判断是否是成功返回的二维码图片
        if (in_array('image/jpeg', $resType)) {
            $image = uploadImage($res, ImageUtil::TYPE_QR_CODE);
        } else {
            throw new Error('100000', '获取小程序二维码失败，微信返回：' . json_encode($res));
        }

        // 返回路径
        return $image->url;
    }

}
