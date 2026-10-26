<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\DistributionGoods;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use App\Models\PartnerEnterpriseSettings;
use Illuminate\Http\Request;

class BatchSettingController extends BaseController
{
    /**
     * 分销商品批量设置
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['distribution_goods_ids', 'promotion_commission_amount', 'service_commission_amount','enterprise_id']);
        $this->validate($request, [
            'distribution_goods_ids' => 'required',
            'promotion_commission_amount' => 'required',
        ], [
            'distribution_goods_ids.required' => '分销商品id不能为空',
            'promotion_commission_amount.required' => '推广佣金比例不能为空',
        ]);
        $params['service_commission_amount'] = $params['service_commission_amount'] ?? 0;

        if ($params['promotion_commission_amount'] + $params['service_commission_amount'] > 10000) {
            notice(1000001);
        }
        //分销合伙人佣金百分比+导购提成百分比+服务佣金百分比<100%
        $detail = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $detail = !empty($detail) ? $detail->toArray() : [];
            //效验二次分销设置
            if(!empty($detail) && $detail['partner_status'] ==  PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE
                && $detail['commission_type'] == PartnerEnterpriseSettings::COMMISSION_TYPE_GOODS)
            {
                if ( $detail['commission_rate'] + $params['promotion_commission_amount'] + $params['service_commission_amount'] > 10000) {
                    notice(1100031);
                }
            }


        DistributionGoods::query()
            ->whereIn('id', $params['distribution_goods_ids'])
            ->update([
                'promotion_commission_amount' => $params['promotion_commission_amount'],
                'service_commission_amount' => $params['service_commission_amount']
            ]);

        return success();
    }
}
