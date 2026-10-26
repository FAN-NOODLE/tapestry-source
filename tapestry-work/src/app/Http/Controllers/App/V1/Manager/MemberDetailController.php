<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;

class MemberDetailController extends BaseController
{
    /**
     * uc人员删除状态
     */
    const MEM_STATUS_DEL = 3;

    /**
     * 员工详情接口
     *
     * @param Request $request
     * @return array|mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {

        $memUid = $request->input('mem_uid', '');

        if (empty($memUid)) {
            error(200011);
        }

        // 获取用户信息
        $memUser = (new Member())->detail($this->qyDomain, $this->appIdentifier, ['memUid' => $memUid]);

        if (!isset($memUser)) {
            error(60103);
        }

        $status = 1;
        // 如果用户是已删除
        if (self::MEM_STATUS_DEL == $memUser['memStatus']) {
            // 查询该员工是否还存在已标记的客户
            $total = MemberUserRelation::where([
                'mem_uid' => $memUid
            ])->count();

            // 如果存在未分配的客户
            if ($total) {
                // 离职未分配
                $status = 2;
            } else {
                // 离职已分配
                $status = 3;
            }
        }

        // 员工角色
        $role = '';
        if (isset($memUser['role']) && !empty($memUser['role'])) {
            $role = $memUser['role']['roleName'];
        }

        // 调用sdk获取会员统计数据
        $result = (new User())->userCountByMemUid(['enterprise_id' => $this->enterpriseId, 'mem_uid' => $memUid]);

        // 数据返回
        return success(['data' => [
            'mem_uid'            => $memUser['memUid'],
            'mem_username'       => $memUser['memUsername'],
            'mem_face'           => $memUser['memFace'],
            'mem_weixin'         => $memUser['memWeixin'],
            'mem_mobile'         => $memUser['memMobile'],
            'mem_role'           => $role,
            'mem_gender'         => $memUser['memGender'],
            'department_name'    => isset($memUser['dpName']) ? $memUser['dpName'][0]['dpName'] : '',
            'mem_subscribe_time' => !empty($memUser['memSubscribeTime']) ? date('Y-m-d', intval($memUser['memSubscribeTime'] / 1000)) : '',
            'customer_num'       => MemberUserRelation::where(['mem_uid' => $memUser['memUid']])->count(),
            'user_num'           => isset($result['data']) ? $result['data']['total'] : 0,
            'status'             => $status
        ]]);
    }

}
