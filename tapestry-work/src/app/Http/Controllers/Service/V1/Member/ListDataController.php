<?php

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MemberService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ListDataController extends BaseController
{

    /**
     * 导购员工列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request): JsonResponse
    {
        // 接收参数
        $params = $request->only([
            'page',
            'size',
            'enterprise_id',
            'worker_status',
            'type',
            'dp_third_id',
            'member_name',
            'worker_no'
        ]);

        $result = MemberService::getMemberList($params);

        return successJson($result);
    }

}
