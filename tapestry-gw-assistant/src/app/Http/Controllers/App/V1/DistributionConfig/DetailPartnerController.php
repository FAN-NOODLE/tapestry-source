<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/23 18:45
 */

namespace App\Http\Controllers\App\V1\DistributionConfig;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Work\DistributionConfig;
use Illuminate\Support\Facades\Log;

class DetailPartnerController  extends BaseController
{
    public function index(DistributionConfig $distribution)
    {
        $params['enterprise_id'] = $this->enterpriseId;
        $params['store_id']      = $this->request->input('store_id',0);
        try {
            $detail = $distribution->partnerDetail($params);
            //添加生成海报公共图片
            $detail['mini_logo']               = $this->loginAdminEnterprise['logo'];
            // 获取商户信息
//            try {
//                $shopDetail = (new MiniProgram())->detail([
//                    'enterprise_id' => $params['enterprise_id']
//                ]);
//            } catch (Error $error) {
//                if ($error->getCode() == 300100) {
//                    notice(657012);
//                }
//            }
//            $detail['merchant_info'] = $shopDetail['data']['name'] ?? '';
            $memberInfo = $this->getCurrentUserInfo($params);
            if (empty($memberInfo)) {
                notice(655001);
            }
            $detail['user_avatar'] = $memberInfo ? $memberInfo['avatar'] : "";
            $detail['user_avatar'] = str_replace('http:', 'https:', $detail['user_avatar']);
            return success($detail);
        } catch (\Exception $e) {
            Log::error("调用接口Work:Distribution->configDetail出错，错误信息：{$e->getMessage()}, 错误code：{$e->getCode()}", $params);
            return error($e->getCode());
        }
    }
}