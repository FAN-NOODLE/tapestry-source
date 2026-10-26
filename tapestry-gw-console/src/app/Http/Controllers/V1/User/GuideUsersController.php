<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class GuideUsersController extends BaseController
{
    const DEF_MEM_ROLE = '导购';

    /**
     * 导购分配员工列表
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['mem_uid', 'status', 'page', 'size']);
        $rules  = [
            'mem_uid' => 'required',
            'status'  => 'required',
        ];

        $message = [
            'mem_uid.required' => '导购[mem_uid]不能为空',
            'status.required'  => '客户[status]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        // 查詢参数组装
        $where = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'page'           => $params['page'] ?? 1,
            'size'           => $params['size'] ?? 20,
            'params'         => [
                'enterprise_id' => $this->enterpriseId,
                'mem_role'      => self::DEF_MEM_ROLE,
                'mem_uid'       => $params['mem_uid'],
                'status'        => $params['status'],
            ],
        ];

        $result = (new Member())->guideUsers($where);

        return success([
            'list'  => $result['data'],
            'page'  => $result['meta']['page'],
            'size'  => $result['meta']['size'],
            'total' => $result['meta']['total'],
        ]);
    }
}
