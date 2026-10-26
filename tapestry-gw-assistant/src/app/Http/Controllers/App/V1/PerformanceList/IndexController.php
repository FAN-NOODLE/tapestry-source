<?php

namespace App\Http\Controllers\App\V1\PerformanceList;

use App\Common\RedisKey;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use App\Micros\MicroScrm\GoMicroScrm;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Auth\Access\AuthorizationException;
use Illuminate\Http\JsonResponse;
use PHPUnit\Util\Json;

class IndexController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function activeCustomers(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $admin_id                = $this->loginAdminEnterprise['admin_id'];

        $memberStores       = $this->getAuthorizedMemberStores($admin_id);
        $params['third_id'] = $this->resolveAuthorizedThirdIdsFromStores($params['third_id'] ?? null, $memberStores);

        $result = GoMicroScrm::performanceListActiveCustomers($params);
        return success($result);
    }

    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function invitation(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $admin_id                = $this->loginAdminEnterprise['admin_id'];

        $memberStores       = $this->getAuthorizedMemberStores($admin_id);
        $params['third_id'] = $this->resolveAuthorizedThirdIdsFromStores($params['third_id'] ?? null, $memberStores);

        $result = GoMicroScrm::performanceListInvitation($params);
        return success($result);
    }

    /**
     * @return JsonResponse
     */
    public function star(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['page']          = 1;
        $params['size']          = 10;
        $params['third_id']      = [];


        $redisKey = RedisKey::PERFORMANCE_LIST_STAR . $this->enterpriseId . ':' . implode('', $params['financial_time']);
        $result   = getOrSetData($redisKey);
        if ($result) {
            return success($result);
        }

        //绩效榜
        $invitation = GoMicroScrm::performanceListStarInvitation($params);

        //顾客榜
        $params['active_type'] = 4;
        $customer              = GoMicroScrm::performanceListStarActiveCustomers($params);

        //活跃榜
        $params['active_type'] = 3;
        $active                = GoMicroScrm::performanceListStarActiveCustomers($params);

        $result = [
            'invitation' => $invitation,
            'customer'   => $customer,
            'active'     => $active
        ];

        getOrSetData($redisKey, $result);

        return success($result);
    }


    /**
     * 业绩榜单
     *
     * @return void
     */
    public function achievement(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['start']         = 0; //0榜单列表 1明星榜单数据
        $admin_id                = $this->loginAdminEnterprise['admin_id'];

        //选择了店铺
        if (isset($params['third_id']) && $params['third_id']) {
            $params['third_id'] = [$params['third_id']];
        } else {
            // 获取企业管理者权限范围数据
            $dataRange          = (new Account())->dataRangeList(['admin_id' => $admin_id]);
            $dpIds              = collect($dataRange)->pluck('object_id')->all();
            $params['third_id'] = $dpIds;
        }

        $result = GoMicroScrm::achievement($params);
        return success($result);
    }

    /**
     * sw顾客榜单
     *
     * @return void
     */
    public function swCustomers(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['start']         = 0; //0榜单列表 1明星榜单数据
        $admin_id                = $this->loginAdminEnterprise['admin_id'];

        //选择了店铺
        if (isset($params['third_id']) && $params['third_id']) {
            $params['third_id'] = [$params['third_id']];
        } else {
            // 获取企业管理者权限范围数据
            $dataRange          = (new Account())->dataRangeList(['admin_id' => $admin_id]);
            $dpIds              = collect($dataRange)->pluck('object_id')->all();
            $params['third_id'] = $dpIds;
        }

        $result = GoMicroScrm::swCustomers($params);

        return success($result);
    }

    /**
     * sw活跃榜单
     *
     * @return void
     */
    public function swActive(): JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['start']         = 0; //0榜单列表 1明星榜单数据
        $admin_id                = $this->loginAdminEnterprise['admin_id'];

        //选择了店铺
        if (isset($params['third_id']) && $params['third_id']) {
            $params['third_id'] = [$params['third_id']];
        } else {
            // 获取企业管理者权限范围数据
            $dataRange          = (new Account())->dataRangeList(['admin_id' => $admin_id]);
            $dpIds              = collect($dataRange)->pluck('object_id')->all();
            $params['third_id'] = $dpIds;
        }

        $result = GoMicroScrm::swActive($params);
        return success($result);
    }

    public function swStar()
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['page']          = 1;
        $params['size']          = 10;
        $params['third_id']      = [];
        $params['start']         = 1;

        $redisKey = RedisKey::PERFORMANCE_LIST_SW_STAR . $this->enterpriseId . ':' . implode('', $params['financial_time']);
        $result   = getOrSetData($redisKey);
        if ($result) {
            return success($result);
        }

        //业绩效榜
        $achievement = GoMicroScrm::achievement($params);
        //顾客榜
        $customer = GoMicroScrm::swCustomers($params);
        //活跃榜
        $active = GoMicroScrm::swActive($params);

        $result = [
            'achievement' => $achievement,
            'customer'   => $customer,
            'active'     => $active
        ];
        getOrSetData($redisKey, $result);

        return success($result);
    }

    public function ksStar()
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['user_type']     = $this->getRoleRange()['member_type'];
        $params['page']          = 1;
        $params['size']          = 10;
        $params['third_id']      = [];
        $params['start']         = 1;

        $redisKey = RedisKey::PERFORMANCE_LIST_KS_STAR . $this->enterpriseId . ':' . implode('', $params['financial_time']);
        $result   = getOrSetData($redisKey);
        if ($result) {
            return success($result);
        }

        //绩效榜
        $achievement = GoMicroScrm::achievement($params);
        //顾客榜
        $params['active_type'] = 4;
        $customer              = GoMicroScrm::performanceListStarActiveCustomers($params);
        //活跃榜
        $params['active_type'] = 3;
        $active                = GoMicroScrm::performanceListStarActiveCustomers($params);

        $result = [
            'achievement' => $achievement,
            'customer'   => $customer,
            'active'     => $active
        ];
        getOrSetData($redisKey, $result);

        return success($result);
    }

    protected function resolveAuthorizedThirdIds($requestedThirdId, array $allowedThirdIds): array
    {
        $allowedThirdIds = array_map('intval', $allowedThirdIds);
        if (!$requestedThirdId) {
            return $allowedThirdIds;
        }

        $requestedThirdId = (int) $requestedThirdId;
        if (!in_array($requestedThirdId, $allowedThirdIds, true)) {
            throw new AuthorizationException('无权访问');
        }

        return [$requestedThirdId];
    }

    protected function getAuthorizedMemberStores(int $adminId): array
    {
        $dataRange = (new Account())->dataRangeList(['admin_id' => $adminId]);
        $thirdIds = collect($dataRange)->pluck('object_id')->all();
        if (empty($thirdIds)) {
            $thirdIds = [0];
        }

        return Common::microScrm()->setApi('micro-scrm/v1/member/store-lists')->post([
            'enterprise_id' => $this->enterpriseId,
            'third_ids'     => $thirdIds,
            'mobile'        => $this->mobile,
            'is_sm'         => $this->isManager ? 1 : 0,
        ]);
    }

    protected function resolveAuthorizedThirdIdsFromStores($requestedThirdId, array $memberStores): array
    {
        return $this->resolveAuthorizedThirdIds($requestedThirdId, array_column($memberStores, 'store_id'));
    }
}
