<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\HandoverService;
use Illuminate\Http\Request;

class MemberTransferController extends BaseController
{
    /**
     * 客户交接接口
     *
     * @param Request         $request
     * @param HandoverService $handoverService
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function index(Request $request, HandoverService $handoverService)
    {
        // 获取参数
        $memUid    = $request->input('mem_uid', '');
        $oldMemUid = $request->input('old_mem_uid', '');
        $ecUserIds = $request->input('ec_userids', []);

        if (empty($memUid)) {
            error(200011);
        }

        if (empty($oldMemUid)) {
            error(200015);
        }

        if (empty($ecUserIds)) {
            error(200014);
        }

        // 转移客户
        $handoverService->handoverCustomer(
            $this->enterpriseId, $this->qyDomain, $oldMemUid, $memUid, $ecUserIds
        );

        // 数据返回
        return success(['data' => []]);
    }
}
