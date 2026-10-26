<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class CustomerTransferController extends BaseController
{
    /**
     * 导购分配员工列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'from_mem_uid',
            'target_mem_uid',
            'transfer_type',
            'union_ids'
        ]);

        $rules = [
            'from_mem_uid'   => 'required',
            'target_mem_uid' => 'required',
            'transfer_type'  => 'required|in:1,2,3',
        ];

        $message = [
            'from_mem_uid.required'   => '当前导购[from_mem_uid]不能为空',
            'target_mem_uid.required' => '目标导购[target_mem_uid]不能为空',
            'transfer_type.required'  => '转交类型[transfer_type]不能为空',
            'transfer_type.in'        => '转交类型[transfer_type]不在范围',
        ];

        Validator::make($params, $rules, $message)->validate();

        // 查詢参数组装
        $params = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'enterprise_id'  => $this->enterpriseId,
                'from_mem_uid'   => $params['from_mem_uid'],
                'target_mem_uid' => $params['target_mem_uid'],
                'transfer_type'  => $params['transfer_type'],
                'union_ids'      => $params['union_ids'] ?? [],
            ],
        ];

        $result = (new Member())->customerTransfer($params);

        return success($result['data']);
    }
}
