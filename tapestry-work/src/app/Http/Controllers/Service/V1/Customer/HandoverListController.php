<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/4
 * Time: 15:58
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\HandoverListService;
use Illuminate\Http\Request;

class HandoverListController extends BaseController
{
    /**
     * 离职员工未分配/已分配客户列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', 0);
        // 查询条件
        $search = $request->get('search', []);

        // 离职员工id
        $memUid = $search['mem_uid'] ?? '';
        if (!$memUid) {
            error(600006);
        }

        // 企业id
        $enterpriseId = $request->get('enterprise_id', 0);
        // 企业domain
        $qyDomain = $request->get('qy_domain', '');
        // 企业应用标识
        $appIdentifier = $request->get('app_identifier', '');
        // 分配状态（1：未分配，2：已分配），默认1
        $handoverStatus = $search['status'] ?? Constant::CUSTOMER_HANDOVER_STATUS_FALSE;

        list($total, $data) = (new HandoverListService())->getHandoverList($enterpriseId, $qyDomain, $appIdentifier,
            $memUid, $handoverStatus, $page, $size);

        return success([
            'data' => $data,
            'meta' => [
                'total' => intval($total),
                'page'  => intval($page),
                'size'  => intval($size)
            ]
        ]);
    }
}
