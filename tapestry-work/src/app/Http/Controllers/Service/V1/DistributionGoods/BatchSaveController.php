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

class BatchSaveController extends BaseController
{
    /**
     * 分销商品批量保存
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['params', 'enterprise_id']);
        $this->validate($request, [
            'enterprise_id' => 'required',
            'params' => 'required',
            'params.*.spu_id' => 'required',
            'params.*.sku_id' => 'required',
            'params.*.distribution_start' => 'required',
            'params.*.distribution_end' => 'required'
        ], [
            'enterprise_id.required' => '企业id不能为空',
            'params.required' => '分销商品不能为空',
            'params.*.spu_id.required' => '分销商品spuid不能为空',
            'params.*.sku_id.required' => '分销商品skuid不能为空',
            'params.*.distribution_start.required' => '分销开始时间不能为空',
            'params.*.distribution_end.required' => '分销结束时间不能为空',
        ]);

        //分销合伙人佣金百分比+导购提成百分比+服务佣金百分比<100%
        $detail = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $detail = !empty($detail) ? $detail->toArray() : [];

        $insert = $params['params'];
        foreach ($insert as &$param) {
            $promotionCommissionAmount = $param['promotion_commission_amount'] ?? 0;
            $serviceCommissionAmount = $param['service_commission_amount'] ?? 0;
            //效验二次分销设置
            if(!empty($detail) && $detail['partner_status'] ==  PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE
                && $detail['commission_type'] == PartnerEnterpriseSettings::COMMISSION_TYPE_GOODS)
            {
                if ( $detail['commission_rate'] + $promotionCommissionAmount + $serviceCommissionAmount > 10000) {
                    notice(1100031);
                }
            }else{//没设置或者选择基于推广佣金分账
                if ($promotionCommissionAmount + $serviceCommissionAmount > 10000) {
                    notice(1000001);
                }
            }


            $param['enterprise_id'] = $params['enterprise_id'];
            $param['promotion_commission_amount'] = $promotionCommissionAmount;
            $param['service_commission_amount'] = $serviceCommissionAmount;
            $param['created_at'] = time();
            $param['updated_at'] = time();
        }
        unset($param);
        DistributionGoods::query()
            ->insert($insert);

        return success();
    }
}
