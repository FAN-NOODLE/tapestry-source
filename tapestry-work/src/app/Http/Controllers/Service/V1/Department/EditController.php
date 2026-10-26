<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:35
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use App\Library\Uc\Member;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class EditController extends BaseController
{
    /**
     * 编辑部门
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        Log::debug('编辑部门参数：'.json_encode($params));

        // 非门店
        if (isset($params['dptId']) && !$params['dptId']) {
            // 部门id
            $dpId = $params['dpId'];
            // 查询部门详情
            $department = (new Department())->detail($qyDomain, $appIdentifier, $dpId);

            Log::debug('编辑部门-部门详情：'.json_encode($department));

            if (isset($department['dptId']) && $department['dptId'] != $params['dptId']) {
                // 查询部门员工列表
                $members = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, ['dpId' => $dpId]);

                Log::debug('编辑部门-员工列表：'.json_encode($members['list']));

                // 员工id集合
                $memUids = [];
                foreach ($members['list'] as $value) {
                    // 只处理有角色的
                    if ($value['memRole']) {
                        $memUids[] = $value['memUid'];
                    }
                }

                // 分组
                $groupMemUids = array_chunk($memUids, 100);

                Log::debug('编辑部门-员工分组：'.json_encode($groupMemUids));

                foreach ($groupMemUids as $memUids) {
                    $memberParams = [
                        'memRole' => '无',
                        'memUids' => $memUids
                    ];
                    (new Member())->updateRole($qyDomain, $appIdentifier, $memberParams);
                }
            }
        }

        $result = (new Department())->edit($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}
