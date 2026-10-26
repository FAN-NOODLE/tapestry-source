<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Demo;
use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\V1\BaseController;

class GuideController extends BaseController
{
    const DEF_MEM_ROLE = '导购';

    /**
     * 导购分配员工列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $page   = $request->input('page', 1);
        $size   = $request->input('size', 20);
        $params = $request->only(['mem_user_name', 'mem_mobile', 'dp_id']);

        // 查詢参数组装
        $guideParams = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'enterprise_id' => $this->enterpriseId,
                'page'          => $page,
                'size'          => $size,
                'mem_user_name' => $params['mem_user_name'] ?? '',
                'mem_mobile'    => $params['mem_mobile'] ?? '',
                'dp_id'         => $params['dp_id'] ?? '',
            ],
        ];

        $result = (new Member())->guideList($guideParams);

        return success(array_merge($result['data'], [
            'page'  => $result['meta']['page'],
            'size'  => $result['meta']['size'],
            'total' => $result['meta']['total'],
        ]));
    }

    public function list()
    {
        // 查詢参数组装
        $params = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'enterprise_id'  => $this->enterpriseId
        ];

        $result = (new Member())->guide($params);

        return success($result);
    }

    public function getDeptList()
    {
        // 查詢参数组装
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'identifier'    => $this->appIdentifier,
            'domain'        => $this->qyDomain
        ];

        $result = (new Demo())->getDeptList($params);

        return success($result);
    }
}
