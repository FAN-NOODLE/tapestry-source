<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StaffShare;
use App\Services\StaffManagerService;
use App\Services\StaffSettingService;
use App\Services\StaffShareService;
use Illuminate\Support\Facades\Log;

class StaffPurchaseDetailController extends BaseController
{
    /**
     * 内购详情接口
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params['open_mem_uid']    = $this->request->input('open_mem_id');
        $params['user_id']         = $this->request->input('user_id');
        $params['enterprise_id']   = $this->request->input('enterprise_id');
        $params['enterprise_hash'] = $this->request->input('enterprise_hash');
        $params['mem_uid']         = $this->request->input('mem_uid');
        $params['parent_user_id']  = $this->request->input('parent_user_id');
        Log::info('时间:' . date('Ymd H:i:s') . ',内购页进入,mem_uid:' . $params['mem_uid'] . ',parent_id:' . $params['parent_user_id'] . ',user_id:' . $params['user_id']);
        // 参数验证
        $this->validateParams($params);

        $staffManagerService = new StaffManagerService();
        $staffSettingService = new StaffSettingService();
        $staffShareService   = new StaffShareService();
        // 管理后台异常实时处理
        $isStaff             = $staffManagerService->staffIsRight($params['open_id'], $params['enterprise_hash']);
        // 当前用户分享数据
        $staffSettingService->handleBindUser($params['user_id'], $isStaff, $params['enterprise_id']);
        $shareUserInfo = (new StaffShare())->getStaffShareInfoByCondition(['user_id' => $params['user_id']], false);
        // 上一级用户数据
        [$params['mem_uid'], $shareParentUserInfo, $params['parent_user_id']] = $staffShareService->getShareParentInfo($params['mem_uid'],
            $params['parent_user_id'], $params['user_id'], $this->enterpriseHash, $shareUserInfo, $isStaff);
        // 抢客逻辑
        if (!empty($shareUserInfo) && $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_YES) {
            $bool = $staffSettingService->changeOwner($params['parent_user_id'], $params['mem_uid'], $shareUserInfo, $shareParentUserInfo, $isStaff);
            // 未抢客 展示原来绑定mem_uid
            if (!$bool) {
                $params['mem_uid'] = $shareUserInfo['mem_uid'];
            }
        }
        // 第一次分享进来
        if (empty($shareUserInfo) || $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO) {
            $params['mem_uid'] = $staffShareService->createStaffShareData($params['mem_uid'], $params['parent_user_id'], $params['user_id'],
                $params['enterprise_id'], $shareParentUserInfo, $shareUserInfo, $isStaff);
        }
        // 获取链路数据
        [$shareList, $shortName] = $staffShareService->getStaffShareList($params['mem_uid'], $params['parent_user_id'], $params['user_id'], $params['enterprise_hash'], $isStaff, $shareUserInfo);
        // 链路等级
        $shareLevel = empty($shareParentUserInfo) ? 0 : ($shareParentUserInfo['level'] == 0 ? 1 : $shareParentUserInfo['level'] + 1);
        // 员工自己进入n初始化
        if (last($shareList)['is_staff']) {
            $shareLevel = 0;
        }
        // 如果解绑其他未绑定普通用户n=2(产品需求)
        if (empty($shareUserInfo) && !empty($shareParentUserInfo) && $shareParentUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO) {
            $shareLevel = 2;
        }

        return success([
            'share_list'  => $shareList,
            'share_level' => $shareLevel,
            'short_name'  => $shortName
        ]);
    }

    public function validateParams($params)
    {
        if ((empty($params['mem_uid']) && empty($params['parent_user_id'])) || (!empty($params['mem_uid']) && !empty($params['parent_user_id']))) {
            notice(720011);
        }
        if (empty($params['open_mem_uid']) || empty($params['user_id']) || empty($params['enterprise_id']) || empty($params['enterprise_hash'])) {
            notice(720011);
        }
    }
}
