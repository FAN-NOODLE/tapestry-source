<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use App\Common\Constant;
use Illuminate\Http\Request;
use CrsSdk\Work\MemberExtend;
use Illuminate\Http\JsonResponse;
use App\Http\Requests\UserRequest;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 新增员工
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        if (empty($this->qyDomain)) {
            return success();
        }

        // 参数列表
        $params = $request->post();

        // 参数验证
        $validator = (new UserRequest())->add($params);
        if ($validator->fails()) {
            error($validator->messages()->first());
        }

        // 部门是否被设置为门店（0:否，1:是）
        $isStore = intval($params['is_store']);
        // 员工角色
        $memRole = $params['mem_role'] ?? '';

        // 不是门店
        if (Constant::DEPARTMENT_IS_STORE_FALSE === $isStore) {
            $memRole = '';
        }

        // 头像
        $memFace = $params['mem_face'] ?? '';

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'corpid'         => $this->enterprise['corp_id'],
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'memUsername' => $params['mem_username'],
                'memFace'     => $memFace,
                'dpIdList'    => [$params['dp_id']],
                'memGender'   => $params['mem_gender'],
                'memMobile'   => $params['mem_mobile'],
                'memRole'     => trim($memRole)
            ]
        ];
        $detail    = (new Member())->add($condition);

        $save = [
            'enterprise_id'   => $this->enterpriseId,
            'guide_id'        => '',
            'mem_uid'         => $detail['memUid'],
            'avatar'          => $memFace,
            'contact_qr_code' => '',
            'config_id'       => '',
            'role'            => trim($memRole),
        ];
        // 保存员工信息
        (new MemberExtend())->saveExtend($save);

        // 写入记录日志
        $this->writeOperatorLog('新增员工：' . $params['mem_username'] . '，角色：' . $memRole);

        return success();
    }
}
