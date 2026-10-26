<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use App\Exceptions\Error;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class OtherGuidesController extends BaseController
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
        $page   = $request->input('page', 1);
        $size   = $request->input('size', 20);
        $memUid = $request->input('mem_uid', '');

        if (empty(trim($memUid))) {
            throw new Error(10000, '导购[mem_uid]参数缺失');
        }

        // 查詢参数组装
        $params = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'enterprise_id' => $this->enterpriseId,
                'mem_role'      => self::DEF_MEM_ROLE,
                'mem_uid'       => $memUid,
                'page'          => $page,
                'size'          => $size,
            ],
        ];

        $result = (new Member())->otherGuideList($params);

        return success([
            'list'  => $result['data'],
            'page'  => $result['meta']['page'],
            'size'  => $result['meta']['size'],
            'total' => $result['meta']['total'],
        ]);
    }
}
