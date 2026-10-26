<?php


namespace App\Services;

use App\Exceptions\Error;
use App\Models\CrsWorkerThirdPartRelation;
use Illuminate\Support\Facades\Log;
use App\Models\PartnerMember;
use App\Models\PartnerEnterpriseMap;

/**
 * 导购认证
 * Class GuideAuthService
 * @package App\Services
 */
class GuideAuthService
{

    /**
     * 导购绑定用户
     * @param $enterpriseId
     * @param $memUid
     * @param $sourceType
     * @param $sourceUserId
     * @throws \Exception
     */
    public function bindUser($enterpriseId, $memUid, $sourceType, $sourceUserId)
    {
        Log::info("GuideAuthService.bindUser params", func_get_args());
        if (CrsWorkerThirdPartRelation::SYNC_SOURCE_CRS == $sourceType) {
            //查询我们库里的合合伙人，如果是合伙人，则进行绑定
            $partner = PartnerEnterpriseMap::query()
                ->where('partner_id', $memUid)
                ->where('enterprise_id', $enterpriseId)
                ->where('user_id', '=', 0)
                ->first();
            if (!empty($partner)) {
                $partner->user_id = $sourceUserId;
                $partner->save();
            } else {
                StoreWorkerService::workerBindUser([
                    'mem_uid' => $memUid,
                    'enterprise_id' => $enterpriseId,
                    'user_id' => $sourceUserId
                ]);
            }
        } else {
            $workers = StoreWorkerService::list([
                'mem_uids'      => [$memUid],
                'enterprise_id' => $enterpriseId,
            ], false);
            if (!empty($workers['list'])) {
                $worker = current($workers['list']);
                $data = [
                    'sync_source' => $sourceType,
                    'enterprise_id' => $enterpriseId,
                    'third_part_user_id' => $sourceUserId,
                    'store_staff_id' => $worker['id'],
                ];
                CrsWorkerThirdPartRelation::query()->create($data);
            } else {
                Log::info("GuideAuthService.bindUser workers is empty");
                throw new Error(110003, "导购信息不存在", 200);
            }

        }
    }

}