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
use App\Models\MemberExtend;
use Illuminate\Http\Request;

class MemberLeaveController extends BaseController
{
    /**
     * 员工离职接口
     *
     * @param Request $request
     *
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

        // 调用uc接口删除用户，即用户离职
        (new Member())->delete($this->qyDomain, $this->appIdentifier, ['memUid' => $memUid]);

        // 数据返回
        return success(['data' => []]);
    }

}
