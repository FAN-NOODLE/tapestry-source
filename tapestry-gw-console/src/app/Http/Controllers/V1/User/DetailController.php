<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use App\Common\Constant;
use Illuminate\Http\Request;
use CrsSdk\Work\MemberExtend;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 员工详情
     * @param Request $request
     * @return JsonResponse|void
     */
    public function index(Request $request)
    {
        // 员工id
        $memUid = $request->get('mem_uid', '');
        if (!$memUid) {
            error(320151);
        }

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => ['memUid' => $memUid]
        ];
        $detail    = (new Member())->detail($condition);
        // 初始化待返数据
        $data = [];
        // 数据格式化
        if (isset($detail['memUid'])) {
//            // 在职状态：在职
//            $memActive = Constant::MEMBER_ACTIVE_TRUE;
//            // 员工状态：已删除
//            if (isset($detail['memStatus']) && Constant::MEMBER_STATUS_DELETED === intval($detail['memStatus'])) {
//                // 在职状态：离职
//                $memActive = Constant::MEMBER_ACTIVE_FALSE;
//            }
//
//            // 批量获取员工头像列表
//            $avatarList = (new MemberExtend())->extendList([
//                'enterprise_id' => $this->enterpriseId,
//                'mem_uids'      => [$memUid]
//            ]);
//
//            $avatarList = array_column($avatarList['data'], null, 'mem_uid');

            $data = [
                'mem_uid'              => $memUid,
                'mem_username'         => $detail['memUsername'],
                'mem_face'             => $detail['memFace'],
                'dp_id'                => $detail['dpName'][0]['dpId'] ?? '',
                'dp_parent_id'         => $detail['dpName'][0]['dpParentid'] ?? '',
                'dp_name'              => $detail['dpName'] ?? [],
                'mem_role'             => $detail['role']['roleName'] ?? '',
                'mem_gender'           => $detail['memGender'],
                'mem_mobile'           => $detail['memMobile'],
                'mem_active'           => $detail['memActive'],
                'mem_subscribe_status' => $detail['memSubscribeStatus'],
                'join_time'            => sprintf('%.0f', ($detail['memSubscribeTime'] ?? 0) / 1000),
                'quit_time'            => sprintf('%.0f', ($detail['memDeleted'] ?? 0) / 1000),
            ];
        }

        return success($data);
    }
}
