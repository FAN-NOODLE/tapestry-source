<?php
/**
 * Created by PhpStorm.
 * Function: 企业员工业绩排行
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\App\V1\StaffOrder;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;


class MyStaffController extends BaseController
{

    //   protected $requireLogin = false;
    /**
     * 我的业绩数字
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $type      = $request->input('type', 0);
        $isAll     = $request->input('is_all', 0);
        $memUid    = $request->input('mem_uid', '');
        $startTime = $request->input('start_time', '');
        $endTime   = $request->input('end_time', '');
        if (null == $this->memUid) {
            $this->memUid = $memUid;
        }

        $search['enterprise_id']  = $this->enterpriseId;
        $search['mem_uid']        = $this->memUid;
        $search['type']           = $type;
        $search['qy_domain']      = $this->qyDomain;
        $search['app_identifier'] = $this->appIdentifier;
        $search['start_time']     = $startTime;
        $search['end_time']       = $endTime;

        $data = (new StaffOrderService())->getMyStaff($search, $isAll);

        $mem_params = [
            'memUids' => [$this->memUid]
        ];
        // 初始化用户相关
        $staffManagerService = (new StaffManagerService());

        $result = $staffManagerService->getStaffListFromUc($this->qyDomain, $this->appIdentifier, $mem_params);

        $ucUserInfo = $result['list'][0] ?? [];

        $userINfo = [
            'mem_uid'      => $ucUserInfo['memUid'] ?? '',
            'mem_mobile'   => $ucUserInfo['memMobile'] ?? '',
            'mem_username' => $ucUserInfo['memUsername'] ?? '',
            'mem_face'     => $ucUserInfo['memFace'] ?? '',
            'dp_name'      => '',
            'dp_id'        => 0,
            'is_author'    => $data['is_author'] ?? 0
        ];

        // 用户详情
        $data['user_info'] = $userINfo;

        return success(['data' => $data]);
    }

}