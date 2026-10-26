<?php

declare(strict_types=1);

namespace App\Http\Controllers\App\V1\DistributionGoods;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Work\{DistributionConfig, MarketingPromotion};
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

/**
 * 分享商品二维码生成
 * Class ShareGoodsController
 * @package App\Http\Controllers\App\V1\DistributionGoods
 */
class GoodsQrCodeController extends BaseController
{
    /**
     * @param MarketingPromotion $promotion
     * @param DistributionConfig $distributionConfig
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(
        MarketingPromotion $promotion,
        DistributionConfig $distributionConfig)
    {
        $rules                   = [
            'spu_id' => "required",
            'sku_id' => "required",
        ];
        $message                 = [
            'spu_id.required' => "商品[spu_id]必填",
            'sku_id.required' => "商品[sku_id]必填",
        ];
        $params                  = $this->validate($this->request, $rules, $message);
        $params['store_id']      = $this->request->input("store_id", 0);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['type']          = 1;

       // $memberInfo = $this->storeWorkerDetail($this->enterpriseId, $params['store_id']);
        $memberInfo = $this->getCurrentUserInfo($params);
        if (empty($memberInfo)) {
            error(655001);
        }
        $params['mem_uid'] = isset($memberInfo['mem_uid']) ? $memberInfo['mem_uid'] : $memberInfo['partner_id'];
        $params['share_uid'] = isset($memberInfo['mem_uid']) ? $memberInfo['mem_uid'] : $memberInfo['partner_id'];
        $params['source'] = self::SHARE_TYPE;
        $params['id']      = $params['sku_id'];
        $data              = $promotion->poster($params);

        //生成海报过慢 将公共图片去掉
        //$data['user_avatar'] = $memberInfo ? $memberInfo['avatar'] : "";
        //$data['user_avatar'] = str_replace('http:', 'https:', $data['user_avatar']);
        // 获取分销数据

        //$distribution                    = $distributionConfig->detail(['enterprise_id' => $this->enterpriseId]);
        //$data['distribution_poster_img'] = $distribution ? $distribution['poster_img'] : "";
        //$data['mini_logo']               = $this->loginAdminEnterprise['logo'];

        // 获取商户信息
        //$data['merchant_info'] = $this->loginAdminEnterprise['name'] ?? '';
        /*$shopDetail = (new MiniProgram())->detail([
            'enterprise_id' => $params['enterprise_id']
        ]);
        $data['merchant_info'] = $shopDetail['data']['name'] ?? '';*/
        return success(compact('data'));
    }
}