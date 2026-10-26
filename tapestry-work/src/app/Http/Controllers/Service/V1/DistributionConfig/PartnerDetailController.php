<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/20 9:22
 */

namespace App\Http\Controllers\Service\V1\DistributionConfig;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Distribution;
use App\Models\PartnerEnterpriseSettings;
use App\Models\StaffSetting;
use App\Services\PartnerService;
use Illuminate\Http\Request;

class PartnerDetailController  extends BaseController
{
    public function index(Request $request)
    {
        $rules = [
            'enterprise_id'   => 'required|integer',
        ];

        $messages = [
            'enterprise_id.*'   => '企业ID不正确',
        ];

        $input                       = $this->validate($request, $rules, $messages);
        // 查询分销配置详情
        $detail = PartnerEnterpriseSettings::query()
            ->leftJoin('staff_setting','partner_enterprise_settings.enterprise_id','=','staff_setting.enterprise_id')
            ->leftJoin('distribution','partner_enterprise_settings.enterprise_id','=','distribution.enterprise_id')
            ->where(
                [
                    'partner_enterprise_settings.enterprise_id'=>$input['enterprise_id'],
                    'staff_setting.deleted_at'=>0
                ])
            ->select(['partner_enterprise_settings.*','staff_setting.fission_type','distribution.poster_img'])
            ->first();
        $detail = !empty($detail) ? $detail->toArray() : [];
        if(empty($detail)){
            $detail =  StaffSetting::query()
                ->where(
                    [
                        'enterprise_id'=>$input['enterprise_id'],
                        'deleted_at'=>0
                    ])
                ->select(['fission_type'])
                ->first();
            $detail = !empty($detail) ? $detail->toArray() : ['fission_type'=>1];
        }

        //添加公告默认值
        if(!isset($detail['distribution_title']) || empty($detail['distribution_title'])){
            $detail['distribution_title'] = '欢迎成为品牌合伙人';
            $detail['distribution_text']  = '通过微信将下方品牌商品分享给周围的好友购买，订单完成后即可获得佣金返利~';
        }
        //默认海报背景图
        if(!isset($detail['poster_img']) || empty($detail['poster_img'])){
            $detail['poster_img'] = 'https://crs-oss.vchangyi.com/store-manage-wxapp/default.png';
        }

        return success($detail);
    }
}