<?php
/*
 * @Author: your name
 * @Date: 2020-02-11 12:02:35
 * @LastEditTime: 2020-02-24 20:36:32
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \shop\src\app\Jobs\PasswordCouponJob.php
 */

namespace App\Http\Controllers\Service\V1\MarketingPromotion;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\PosterService;
use CrsSdk\Boss\Auth;
use CrsSdk\Shop\AppCode;
use Illuminate\Support\Facades\Log;

class PosterController extends BaseController
{
    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        set_time_limit(0);
        ini_set('memory_limit', -1);

        $type        = Constant::MARKETING_COUPON;
        $operateType = Constant::OPERATE_TYPE_DOWNLOAD;

        $params = request()->all();

        if (isset($params['type']) && $params['type'] == Constant::MARKETING_GOODS) {
            $type = Constant::MARKETING_GOODS;
        }

        switchAppHost(['enterprise_id' => $params['enterprise_id']]);

        try {
            $service = new PosterService();

            // 生成太阳码
            list($path, $qrCodeParams) = $service->getQrCodeParams($type, $params);

            $params['sun_code'] = (new AppCode())->create(['enterprise_id' => $params['enterprise_id'], 'params' => $qrCodeParams, 'path' => $path])['data']['code_url'] ?? '';
            if (empty($params['sun_code'])) {
                return success([]);
            }

            $miniApp = (new Auth())->authDetail(['auth_type' => Constant::APP_TYPE, 'enterprise_id' => $params['enterprise_id']]);
            if (empty($miniApp)) {
                return success([]);
            }

            // 小程序信息
            $appinfo = [
                'name'      => $miniApp['data']['name'],
                'logo'      => $miniApp['data']['logo'],
                'style_rgb' => ''
            ];

            $mediaId = 0;
            if ($type == Constant::MARKETING_COUPON) {
                $mediaId = $service->generateCouponPoster($appinfo, $params, $operateType, '');
            }
            return success(['media_id' => $mediaId, 'sun_code' => $params['sun_code']]);

        } catch (\Exception $e) {
            Log::info('营销推广海报生成失败-1' . $e);

            throw new Error(11000, '海报生成失败,请重试:' . $e->getMessage());
        }

    }
}
