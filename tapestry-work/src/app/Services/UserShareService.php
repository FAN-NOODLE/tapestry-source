<?php

namespace App\Services;

use App\Common\Constant;
use App\Models\PartnerMember;
use App\Models\StaffShare;
use App\Models\StoreWorker;
use App\Services\ShareChain\StaffShareStrategy;
use App\Services\ShareChain\Substance;
use CrsSdk\Work\StaffRight;
use Illuminate\Support\Facades\Log;

/**
 * 设置分享链路关系
 */
class UserShareService extends BaseService
{

    // 企业ID
    protected $sceneEnterpriseId;
    // 门店ID
    protected $sceneStoreId = 0;
    // 用户ID
    protected $loginUserId = 0;
    // 登陆者用户身份 (1:普通用户 2:导购 10:合伙人)
    protected $userType = Constant::LOGIN_IDENTITY_USER;
    // 登陆者worker
    protected $userWorkerId = 0;
    // 登陆者memUid
    protected $userMemUid = '';
    // 分享者用户Id
    protected $sceneUserId = 0;
    // 分享者用户Id
    protected $sceneMemUid = '';
    // 分享者用户Id
    protected $sceneWorkerId = 0;
    // 分享者用户Id
    protected $memberUserId = 0;
    // 分享者类型 (1:普通用户 2:导购 10:合伙人)
    protected $sceneUserType = Constant::LOGIN_IDENTITY_USER;
    // 用户当前的workId
    protected $loginWorkId = 0;
    // 用户当前的workId
    protected $relationLevel = 0;

    /**
     * UserShareService constructor.
     *
     * @param int $enterpriseId
     * @param int $storeId
     * @param int $userId
     * @param string $memberUserId
     * @param int $sceneUserId
     * @param int $workerId
     * @param int $sceneUserType
     */
    public function __construct(int $enterpriseId,
                                int $storeId,
                                int $userId,
                                string $memberUserId,
                                int $sceneUserId,
                                int $workerId,
                                int $sceneUserType
    )
    {

        $this->sceneEnterpriseId = $enterpriseId;
        $this->sceneStoreId      = $storeId;
        $this->loginUserId       = $userId;
        $this->memberUserId      = $memberUserId;
        $this->sceneUserId       = $sceneUserId;
        $this->sceneWorkerId     = $workerId;
        $this->sceneUserType     = $sceneUserType;
    }

    /**
     * 用户分享链路
     *
     * @return bool
     * @throws \Exception
     */
    public function setUserShare()
    {
        // 更变分享者UserId逻辑
//        if ($this->updateShareUserId()) {

        // 设置用户和分享者身份
        $this->setUserIdentity();
        // 记录用户和分享者链路
        $this->insertUserShareRelation();
//        }

        return $this->loginWorkId;
    }

    /**
     * 根据企业导购业绩设置信息判断是否删除或保存用户推广人缓存
     *
     * @return bool
     * @throws \Exception
     */
    private function checkEnterpriseStaffSetting()
    {
        // 获取企业是否设置或是否开启（导购业绩/员工购买商品是否有佣金）
        /*$guideSettingKey = redisKey('ENTERPRISE_GUIDE_SETTING', ['enterprise_id' => $this->sceneEnterpriseId, 'type' => Constant::LOGIN_IDENTITY_GUIDE]);
        $guideSetting    = getOrSetData($guideSettingKey);

        if (empty($guideSetting)) {
            $guideSetting = (new StaffRight())->staffSetting([
                    'enterprise_id' => $this->sceneEnterpriseId,
                    'type'          => Constant::LOGIN_IDENTITY_GUIDE,
                ])['data'] ?? [];

            if (!empty($guideSetting)) {
                getOrSetData($guideSettingKey, $guideSetting);
            }
        }*/

        $guideSetting = (new StaffRight())->staffSetting([
                'enterprise_id' => $this->sceneEnterpriseId,
                'type'          => Constant::LOGIN_IDENTITY_GUIDE,
            ])['data'] ?? [];

        if (!empty($guideSetting)) {
            // 判断导购（员工购买商品有佣金 || 客户分享统计业绩）
            if ($guideSetting['commission_type'] == Constant::GUIDE_COMMISSION_TYPE_YES ||
                $guideSetting['fission_type'] == Constant::GUIDE_FISSION_TYPE ||
                $this->relationLevel == 1
            ) {
                $this->settingUserWorkId();
            } else {
                $this->delUserWorkCache();
            }

        } else {
            // 如果企业未设置，默认绑定用户永久有效
            $this->settingUserWorkId();
        }
        return true;
    }

    /**
     * 清楚用户推广人缓存信息
     */
    private function delUserWorkCache()
    {
        $key = redisKey('USER_SHARE_WORKER_ID', ['user_id' => $this->loginUserId]);
        delRedis($key);
    }

    /**
     * 设置用户分享者workId到缓存
     *
     * @return bool
     */
    private function settingUserWorkId()
    {
        if (!empty($this->loginWorkId)) {
            $key = redisKey('USER_SHARE_WORKER_ID', ['user_id' => $this->loginUserId]);
            getOrSetData($key, $this->loginWorkId, null);
        }
        return true;
    }

    /**
     * 接收分享者信息更变分享者
     *
     * @return bool
     */
    private function updateShareUserId()
    {
        $sceneUserId = 0;

        // 判断是否存在分享者信息
        if (!empty($this->sceneUserId)) {
            $sceneUserId = $this->sceneUserId;
        } elseif (!empty($this->memberUserId)) {

            // 如果存在memberUserId和storeId查询查询门店版本
            if (!empty($this->sceneStoreId)) {
                $where['store_id'] = $this->sceneStoreId;
            }

            $where['type'] = [1, 2, 3];
            $where['mem_uid']       = $this->memberUserId;
            $where['worker_status'] = StoreWorker::IN_SERVICE;
            $storeWorker = StoreWorkerService::detail($where);
            $this->sceneUserId = $storeWorker['user_id'] ?? 0;

            if (empty($storeWorker)) {
                return false;
            } else {
                return true;
            }
        }

        $this->sceneUserId = $sceneUserId ?? 0;
        if (empty($this->sceneUserId)) {
            return false;
        }
        return true;
    }

    /**
     * 获取当前用户导购Id
     */
    private function getUserWorkId()
    {
        // 获取用户当前绑定导购信息
        $shareRelation = StaffShare::query()
            ->where([
                'enterprise_id' => $this->sceneEnterpriseId,
                'user_id'       => $this->loginUserId,
                'type'          => Constant::LOGIN_IDENTITY_GUIDE,
            ])
            ->whereIn('binding_type', [StaffShare::BINDING_TYPE_YES, StaffShare::BINDING_TYPE_PARENT_UNBINDING])
            ->orderBy('created_at','desc')
            ->first();

        if (!empty($shareRelation->worker_id)) {
            $this->loginWorkId   = $shareRelation->worker_id;
            $this->relationLevel = $shareRelation->level;
        }

        return true;
    }

    /**
     * 记录用户分享链路
     *
     * @return array|bool
     */
    private function insertUserShareRelation()
    {
        $staffShareStrategy = new Substance(new StaffShareStrategy());

        // 如果用户是导购身份或合伙人
        if ($this->userWorkerId > 0) {
            Log::info('用户是导购初始化staffInitial传参' . json_encode([
                    'enterpriseId' => $this->sceneEnterpriseId,
                    'userId'       => $this->loginUserId,
                    'userMemUid'   => $this->userMemUid,
                    'userWorkerId' => $this->userWorkerId,
                    'type'         => $this->userType,
                ]));
            $staffShareStrategy->staffInitial($this->sceneEnterpriseId, $this->loginUserId, $this->userMemUid, $this->userWorkerId, $this->userType);
            return true;
        }

        // 如果分享者是导购身份或合伙人
        if ($this->sceneWorkerId > 0) {
            Log::info('分享者是导购初始化staffInitial传参' . json_encode([
                    'enterpriseId'  => $this->sceneEnterpriseId,
                    'sceneUserId'   => $this->sceneUserId,
                    'sceneMemUid'   => $this->memberUserId,
                    'sceneWorkerId' => $this->sceneWorkerId,
                    'type'          => $this->sceneUserType,
                ]));
            $staffShareStrategy->staffInitial($this->sceneEnterpriseId, $this->sceneUserId, $this->memberUserId ?? '', $this->sceneWorkerId, $this->sceneUserType);
        }

        Log::info('调用分享绑定逻辑staffBinding传参' . json_encode([
                'enterpriseId'  => $this->sceneEnterpriseId,
                'userId'        => $this->loginUserId,
                'sceneUserId'   => $this->sceneUserId,
                'sceneMemUid'   => $this->memberUserId,
                'sceneWorkerId' => $this->sceneWorkerId,
                'type'          => $this->sceneUserType,
            ]));
        // 具体的绑定逻辑，包含抢客等判断
        $staffShareStrategy->staffBinding($this->sceneEnterpriseId, $this->loginUserId, $this->sceneUserId, $this->sceneWorkerId, $this->sceneUserType);

        return true;
    }

    /**
     * 获取用户基础信息
     */
    private function setUserIdentity()
    {
        // 获取登录者的合伙人身份
        $loginPartnerInfo = (new PartnerMemberService())->getPartnerInfoByUserId($this->sceneEnterpriseId, $this->loginUserId, PartnerMember::PARTNER_STATUS_ENABLED);
        if (!empty($loginPartnerInfo)) {
            $this->userType = Constant::LOGIN_IDENTITY_PARTNER;
            $this->userWorkerId = $loginPartnerInfo['id'];
        } else {
            // 获取登录者的导购身份
            $loginWorkerInfo = StoreWorkerService::detail([
                'enterprise_id' => $this->sceneEnterpriseId,
                'user_id' => $this->loginUserId,
                'worker_status' => StoreWorker::IN_SERVICE
            ]);

            if (!empty($loginWorkerInfo)) {
                $this->userType     = Constant::LOGIN_IDENTITY_GUIDE;
                $this->userMemUid   = $loginWorkerInfo['mem_uid'];
                $this->userWorkerId = $loginWorkerInfo['id'];
            }
        }

        return true;
    }

}
