<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/20 9:22
 */

namespace App\Http\Controllers\Service\V1\DistributionConfig;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseSettings;
use App\Models\StaffSetting;
use App\Models\WithdrawalSetting;
use App\Services\PartnerService;
use Illuminate\Http\Request;

class PartnerSaveController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->all();
        //先查看导购提现和导购业绩设置
        $staffSetting = StaffSetting::query()
            ->where('enterprise_id',$input['enterprise_id'])
            ->first();
        $withdrawal = WithdrawalSetting::query()
            ->where('enterprise_id',$input['enterprise_id'])
            ->first();

        if(!empty($staffSetting) && $staffSetting->fission_type != StaffSetting::FISSION_TYPE_TRUE){
            notice('1100027');
        }

        if(empty($staffSetting) && empty($withdrawal)){
            notice('1100026');
        }

        if(empty($staffSetting) && !empty($withdrawal)){
            notice('1100024');
        }
        if(!empty($staffSetting) && empty($withdrawal)){
            notice('1100025');
        }

        $input = $this->getParams($request);

        $params = [
            'enterprise_id' => $input['enterprise_id'],
            'commission_type' => $input['commission_type'],
            'commission_rate' => $input['commission_rate'],
            'withdraw_status' => $input['withdraw_status'],
            'partner_status' => $input['partner_status'],
            'opt_admin_id' => $input['opt_admin_id'] ?? 0,
            'opt_admin_name' => $input['opt_admin_name'] ?? '',
            'distribution_title' => $input['distribution_title'],
            'distribution_text'  => $input['distribution_text'],
            // 'id' => $input['id'] ?? 0
        ];
        try {
            //设置二级分销
             (new PartnerService())->check($params);
            (new PartnerEnterpriseSettings())->insertOrUpdate($params);

        } catch (Error $error) {
            \error($error->getCode(), $error->getMessage());
        } catch (\Exception $exception) {
            notice('partner_enterprise_settings_save_error');
        }
        return success();
    }

    /**
     * 获取请求参数
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams(Request $request): array
    {
        $rules = [
            'partner_status' => 'required|in:0,1',
            'commission_type' => 'required|in:1,2',
            'commission_rate' => 'required',
            'withdraw_status' => 'required|in:0,1',
            'enterprise_id' => 'required',
            'opt_admin_id'=>'integer',
            'opt_admin_name'=>'string',
            'distribution_title'=> 'required|string',
            'distribution_text' => 'required|string'
        ];

        $message = [
            'partner_status.*' => '是否开启分销不合法',
            'commission_type.*' => '分佣方式不合法',
            'commission_rate.*' => '佣金比例不合法',
            'withdraw_status.*' => '是否可提现不合法',
            'enterprise_id.*' => '企业ID不合法',
            'distribution_title.*'=>'公告标题不合法',
            'distribution_text.*'=>'公告内容不合法'
        ];
        $input = $this->validate($request, $rules, $message);
        return $input;
    }
}