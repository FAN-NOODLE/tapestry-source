<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\HandoverListService;
use Illuminate\Http\Request;

class UserListController extends BaseController
{
    /**
     * 导购名下的客户列表(已分配、未分配)接口
     *
     * @param Request             $request
     * @param HandoverListService $handoverListService
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request, HandoverListService $handoverListService)
    {
        // 获取导购uid
        $memUid = $request->input('mem_uid', '');

        // 获取未分配的客户列表
        $result = $handoverListService->getHandoverList($this->enterpriseId, $this->qyDomain, $this->appIdentifier, $memUid, 1);

        // 获取已分配的客户列表
        $resultChange = $handoverListService->getHandoverList($this->enterpriseId, $this->qyDomain, $this->appIdentifier, $memUid, 2);

        // 返回数据
        return success([
            'data' => [
                'list'        => $result[1],
                'change_list' => $resultChange[1]
            ]
        ]);
    }
}
