<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/20 9:40
 */

namespace App\Services;


use App\Exceptions\Error;
use App\Models\DistributionGoods;
use App\Models\PartnerEnterpriseSettings;

class PartnerService
{

    /**
     *
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/20 10:39
     */
    public function check($input)
    {
        //是否开启二级分销 是否基于分销商品
        if($input['partner_status'] == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE
         && $input['commission_type'] == PartnerEnterpriseSettings::COMMISSION_TYPE_GOODS)
        {
            //分销商品列表 进行中的/即将开始的
            $disGoodsList =  DistributionGoods::query()
                ->where('enterprise_id', $input['enterprise_id'])
               //结束时间小于等于当前时间（包含即将开始和进行中的）
                ->Where('distribution_end','>=',time())
                ->get()
                ->toArray();
            //批量效验 分销合伙人佣金百分比+导购提成百分比+服务佣金百分比<100%
            foreach ($disGoodsList as $k=>$v){
                $rate = $v['promotion_commission_amount']+$v['service_commission_amount']+$input['commission_rate'];
                if($rate > 10000) {//
                     notice(1100021);
                    }
            }
        }
        return true;
    }
}