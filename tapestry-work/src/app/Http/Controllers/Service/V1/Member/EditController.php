<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Services\MemberService;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class EditController extends BaseController
{
    /**
     * 编辑员工
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
        $enterpriseId  = $request->post('enterprise_id', 0);
        $dpId          = $request->post('dp_id', '');
        $params        = $request->post('params', []);

        if (empty($params['memFace'])) {
            throw new Error(100000, '员工头像不能为空!');
        }

        if (empty($params['memRole'])) {
            $params['memRole'] = '无';
        }

        try {
            // 更新员工员工内购分享绑定关系确定或解除
            (new StaffManagerService())->editMemberUnbindStaffShare($qyDomain, $appIdentifier, $enterpriseId, $dpId, $params);
        } catch (\Exception $exception) {
            Log::error('更新员工员工内购分享绑定关系确定或解除异常：' . $exception->getMessage());
        }
        
        $params['avatarMediaid'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $params['memFace']);

        unset($params['memFace']);

        $result = (new Member())->edit($qyDomain, $appIdentifier, $params);

        try {
            // 更新员工员工内购部门
            (new StaffManagerService())->editMemberStaffRight($qyDomain, $appIdentifier, $enterpriseId, $dpId, $params);
        } catch (\Exception $exception) {
            Log::error('更新员工员工内购部门异常：' . $exception->getMessage());
        }

        return success($result);
    }
}
