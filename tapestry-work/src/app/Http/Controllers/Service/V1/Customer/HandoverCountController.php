<?php
/**
 * Created by PhpStorm.
 * User: meijie
 * Date: 2019-03-23
 * Time: 23:24
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Services\HandoverCountService;
use Illuminate\Http\Request;

class HandoverCountController
{
    /**
     * 获取离职员工集合交接数据 - 待确认/待分配数量
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $memUids = $request->post('mem_uids', []);
        $enterpriseId = $request->post('enterprise_id', 0);
        if (!$memUids) {
            error(600007);
        }
        // 企业id
        $data = (new HandoverCountService())->getHandoverCount($enterpriseId, $memUids);
        return success(['data' => $data]);
    }
}
