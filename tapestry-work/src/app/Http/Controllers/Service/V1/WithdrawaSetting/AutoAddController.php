<?php
/**
 * 提现配置
 */


namespace App\Http\Controllers\Service\V1\WithdrawaSetting;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;
use App\Models\WithdrawalApplication;
use App\Models\WithdrawalAutomaticReview;
use Illuminate\Http\Request;

class AutoAddController extends BaseController
{


    public function index(Request $request)
    {
        $rules = [
            'enterprise_id' => 'required',
            'auto_status' => 'required|in:1,0',
            'auto_at' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'auto_status.required' => '是否开启状态不能为空',
            'auto_at.required' => '自动批准时间不能为空',
            'auto_status.in'   => '是否开启状态值必须为0或1',
        ];

        $this->validate($this->request, $rules, $messages);

        $memType = $request->get('mem_type', WithdrawalApplication::APP_TYPE_STAFF);
        if (!in_array($memType, [WithdrawalApplication::APP_TYPE_STAFF, WithdrawalApplication::APP_TYPE_PARTNER])) {
            $memType = WithdrawalApplication::APP_TYPE_STAFF;
        }
        $request = $request->input();

        $setInfo = WithdrawalAutomaticReview::query()
            ->where('enterprise_id', $request['enterprise_id'])
            ->where('mem_type', $memType)
            ->first();
        $data = [
            'enterprise_id' => $request['enterprise_id'],
            'auto_status' => $request['auto_status'],
            'auto_at' => $request['auto_at'],
            'mem_type' => $memType,
            'updated_at' => time()
        ];
        if($request['auto_status'] == 1){
            $data['last_at'] = time();
        }
        if(empty($setInfo)){
            $data['created_at'] = time();
            WithdrawalAutomaticReview::query()->insert($data);
        }else{
            WithdrawalAutomaticReview::query()
                ->where('enterprise_id', $request['enterprise_id'])
                ->where('mem_type', $memType)
                ->update($data);
        }
        return success([]);

    }
}