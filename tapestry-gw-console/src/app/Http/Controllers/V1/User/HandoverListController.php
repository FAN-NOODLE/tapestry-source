<?php

namespace App\Http\Controllers\V1\User;

use App\Common\Constant;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class HandoverListController extends BaseController
{
    /**
     * 离职员工未分配/已分配客户列表
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);
        // 查询条件
        $search = $request->get('search', []);

        // 离职员工id
        $memUid = $search['mem_uid'] ?? '';
        if (!$memUid) {
            error(320151);
        }

        // 调用sdk 获取交接客户列表
        $params = [
            'enterprise_id'  => $this->enterpriseId,
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'page'           => $page,
            'size'           => $size,
            'search'         => $search
        ];
        $result = (new Customer())->handoverList($params);

        return success([
            'list'  => $result['data'],
            'total' => intval($result['meta']['total']),
            'page'  => intval($page),
            'size'  => intval($size)
        ]);
    }
}
