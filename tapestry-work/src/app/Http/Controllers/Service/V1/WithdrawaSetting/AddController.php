<?php
/**
 * 提现配置
 */


namespace App\Http\Controllers\Service\V1\WithdrawaSetting;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseSettings;
use App\Models\WithdrawalSetting;
use Illuminate\Http\Request;

class AddController extends BaseController
{


    public function index(Request $request)
    {
        $rules = [
            'enterprise_id' => 'required',
            'transfer_no' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'type.required' => '可提现人员id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        $setInfo = WithdrawalSetting::query()->where('enterprise_id', $request['enterprise_id'])->first();
        $data = [
            'enterprise_id' => $request['enterprise_id'],
            'transfer_no' => $request['transfer_no'],
            'type' => implode($request['type'],','),
            'updated_at' => time()
        ];
        if(empty($setInfo)){
            $data['created_at'] = time();
            WithdrawalSetting::query()->create($data);
        }else{
            WithdrawalSetting::query()->where('enterprise_id', $request['enterprise_id'])->update($data);
        }
        //统一二级分销合伙人提现设置
        $setting = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $request['enterprise_id'])
            ->first();
        if(!empty($setting)){
            //开启合伙人提现
            if(in_array(WithdrawalSetting::PARTNER_STATUS,$request['type'])){
                $setting->update(['withdraw_status'=>PartnerEnterpriseSettings::WITHDRAW_STATUS_TRUE]);
            }
            if(!in_array(WithdrawalSetting::PARTNER_STATUS,$request['type'])){
                $setting->update(['withdraw_status'=>PartnerEnterpriseSettings::WITHDRAW_STATUS_FALSE]);
            }
        }

        return success([]);

    }
}