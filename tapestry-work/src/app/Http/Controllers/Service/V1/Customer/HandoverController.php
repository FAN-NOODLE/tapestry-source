<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/4
 * Time: 19:14
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\HandoverService;
use Illuminate\Http\Request;

class HandoverController extends BaseController
{
    /**
     * 交接客户
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->post('enterprise_id', 0);
        // 企业domain
        $qyDomain = $request->post('qy_domain', '');
        // 离职员工id
        $memUid = $request->post('mem_uid', '');
        // 新导购id
        $newMemUid = $request->post('new_mem_uid', '');
        // 待分配外部联系人id集合
        $ecUserIds = $request->post('ec_userids', []);

        (new HandoverService())->handoverCustomer($enterpriseId, $qyDomain, $memUid, $newMemUid, $ecUserIds);

        return success(['data' => []]);
    }
}
