<?php
/*
 * @Author: your name
 * @Date: 2020-02-11 12:02:35
 * @LastEditTime: 2020-02-24 20:36:32
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \shop\src\app\Jobs\PasswordCouponJob.php
 */

namespace App\Http\Controllers\Html\V1\MarketingPromotion;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Html\V1\BaseController;
use App\Services\CommonService;
use App\Services\PosterService;
use CrsSdk\Boss\Auth;
use CrsSdk\Shop\AppCode;
use CrsSdk\Shop\Template;
use Illuminate\Support\Facades\Log;

class PosterController extends BaseController
{
    /**
     * @return array
     */
    public function index()
    {
        set_time_limit(0);
        ini_set('memory_limit', -1);

        $type = request()->input('type') ?? Constant::MARKETING_GOODS;
        $params = request()->input('params') ?? [];
        $params['mem_uid'] = $this->memUid;
        $params['enterprise_id'] = $this->enterpriseId;
        $operateType = request()->input('operate_type') ?? Constant::OPERATE_TYPE_SEND;

        switchAppHostByDomain($this->domain);

        try{
            $service = new PosterService();

            // 生成太阳码
            list($path, $qrCodeParams) = $service->getQrCodeParams($type, $params);

            // 防止键名过长
            $key = md5(json_encode($qrCodeParams));
            $key = redisKey('QR_CODE', ['key' => $key]);

            $qrCode = getOrSetData($key);
            if (!empty($qrCode)) {
                $params['sun_code'] = $qrCode['sun_code'];
            } else {
                $params['sun_code'] = (new AppCode())->create(['enterprise_id' => $this->enterpriseId, 'params' => $qrCodeParams, 'path' => $path])['data']['code_url'] ?? '';
                if (empty($params['sun_code'])) {
                    return success([]);
                }

                getOrSetData($key, ['sun_code' => $params['sun_code']], 86400);
            }

            // 获取员工信息
            $member = $this->member;

            if ($type == Constant::MARKETING_GOODS) {
                // 商品海报
                $mediaId = $service->generateGoodsPoster($params, $member, $operateType, $this->domain);
                return success(['media_id' => $mediaId]);
            } else {
                $miniApp = (new CommonService())->getAuthDetail($this->enterprise['hash']);
                if (empty($miniApp)) {
                    return success([]);
                }

                // 小程序信息
                $appinfo = [
                    'name' => $miniApp['name'],
                    'logo' => $miniApp['logo'],
                    'style_rgb' => ''
                ];

                if ($type == Constant::MARKETING_COUPON) {
                    // 优惠券海报
                    $mediaId = $service->generateCouponPoster($appinfo, $params, $operateType, $this->domain);
                    return success(['media_id' => $mediaId]);
                } else {
                    $colorKey = redisKey('SHOP_COLOR', ['enterprise_id' => $this->enterpriseId]);

                    $style = getOrSetData($colorKey);
                    if (empty($style)) {
                        // 获取商城主色调
                        $style = (new Template())->shopStyleInfo(['enterprise_id' => $this->enterpriseId])['data'];
                        if (empty($style)) {
                            return success([]);
                        }

                        getOrSetData($colorKey, $style, 7200);
                    }

                    $appinfo['style_rgb'] = substr(explode('),rgb(', $style['rgb_colors'])[0], 4);
                    $appinfo['style_rgb'] = explode(',', $appinfo['style_rgb']);

                    // 活动海报
                    $mediaId = $service->generateActivityPoster($appinfo, $params, $member, $operateType, $this->domain);
                    return success(['media_id' => $mediaId]);
                }
            }

        } catch(\Exception $e) {
            Log::info('营销推广海报生成失败-2' . $e);

            throw new Error(11000, '海报生成失败,请重试');
        }

    }
}
