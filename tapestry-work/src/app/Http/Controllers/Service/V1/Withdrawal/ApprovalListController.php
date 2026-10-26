<?php
/**
 * 提现审批列表
 * Created by PhpStorm
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\WithdrawalApplication;
use App\Services\WithdrawalApplicationService;

class ApprovalListController extends BaseController
{
    /**
     * 提现审批列表
     **
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 接值
        $input = $this->request->only([
            'page',
            'size',
            'enterprise_id',
            'ids',
            'search_type',
            'search_value',
            'start_time',
            'end_time',
            'status',
            'worker_ids'
        ]);
        $rules = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        $input['mem_type'] = $this->request->get('mem_type', 0);
        if (!in_array(
            $input['mem_type'],
            [WithdrawalApplication::APP_TYPE_STAFF, WithdrawalApplication::APP_TYPE_PARTNER]
        )
        ) {
            $input['mem_type'] = WithdrawalApplication::APP_TYPE_STAFF;
        }

        $data = (new WithdrawalApplicationService())->getList($input, $input['is_page']??true);

        return success($data);

    }
}
