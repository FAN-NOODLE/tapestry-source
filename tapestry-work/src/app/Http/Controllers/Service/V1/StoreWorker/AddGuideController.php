<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class AddGuideController extends BaseController
{

    public function index()
    {
        $rules = [
            'enterprise_id'    => 'required',
            'member'           => 'required',
            'operator_id'      => 'required',
            'operator_name'    => 'required',
            'type'             => 'required',
            'range'            => 'array',
            'range_type'       => 'integer',
        ];
        $message = [
            'enterprise_id.required' => '企业id参数必填',
            'operator_id.required'   => '操作人id参数必填',
            'operator_name.required' => '操作人名称参数必填',
            'member.required'        => '导购信息参数必填',
            'type.required'          => '导购类型参数必填',
        ];
        $input = $this->validate($this->request, $rules, $message);
        // 企业管理角色
        if ($input['type'] == Constant::STORE_WORKER_TYPE_MANAGER) {
            if (empty($input['range']) || !isset($input['range_type'])) {
                notice('range_not_empty');
            }
            $range = [];
            foreach ($input['range'] as $item) {
                $range[] = [
                    'type' => $input['range_type'],
                    'object_id' => $item
                ];
            }
            $input['range'] = $range;
        }else{
            $input['type'] = Constant::STORE_WORKER_TYPE_ENTERPRISE;
        }
        // 手机号为空验证
        if (is_array($input['member'])) {
            foreach ($input['member'] as $member) {
                if (empty($member['mobile'])) {
                    notice(2300015, [$member['user_name']]);
                }
            }
        }

        StoreWorkerService::addGuide($input);
        return successful();
    }

}
