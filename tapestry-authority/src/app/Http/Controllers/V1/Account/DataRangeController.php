<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Micros\StoreMicro;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminService;
use Illuminate\Http\Request;

class DataRangeController extends BaseController
{
    public function index(Request $request)
    {
        $adminId = trim($request->post('admin_id'));
        $client  = trim($request->post('client', 'console'));

        if (empty($adminId)) {
            notice(700001);
        }

        $accountDetail = (new EnterpriseAdminService())->enterpriseAdminDetail($adminId, $client);

        $rangeType    = $accountDetail['range_type'];
        $enterpriseId = $accountDetail['enterprise_id'];


        $dataRangeList = (new EnterpriseDataRange())->accountDataRange($adminId);

        // list($storeIds, $hasShopAuth) = $this->getAllStores($rangeType, $dataRangeList, $enterpriseId);

        return success([
            'data_range_list' => $dataRangeList,
            'store_ids'       => [],
            'has_shop_auth'   => 2
        ]);
    }

    protected function getAllStores($rangeType, $dataRangeList, $enterpriseId)
    {
        $allStoreIds = [];
        $storeIds    = [];
        $districtIds = [];
        $hasShopAuth = 2;

        if ($rangeType == EnterpriseAdmin::RANGE_TYPE_ALL) {
            $allStoreIds = (new StoreMicro())->getStoreIds(['enterprise_id' => $enterpriseId]);
            return [$allStoreIds, 1];
        }

        foreach ($dataRangeList as $v) {
            $type     = $v['type'];
            $objectId = $v['object_id'];

            switch ($type) {
                case EnterpriseDataRange::TYPE_ALL:
                    $allStoreIds = (new StoreMicro())->getStoreIds(['enterprise_id' => $enterpriseId]);
                    break 2;
                case EnterpriseDataRange::TYPE_AREA:
                    $districtIds[] = $objectId;
                    break;
                case EnterpriseDataRange::TYPE_STORE:
                    $storeIds[] = $objectId;
                    break;
                case EnterpriseDataRange::TYPE_SHOP:
                    $hasShopAuth = 1;
                    $storeIds[]  = $enterpriseId;
            }

            unset($type);
            unset($objectId);
        }

        if (empty($allStoreIds) && !empty($districtIds) && !empty($enterpriseId)) {
            $allStoreIds = (new StoreMicro())->getStoreIds(['enterprise_id' => $enterpriseId, 'area_ids' => $districtIds]);
        }

        if (!empty($storeIds)) {
            $allStoreIds = array_values(array_unique(array_merge($storeIds, $allStoreIds)));
        }

        return [$allStoreIds, $hasShopAuth];
    }

    /**
     * 获取管理员管理范围
     * @return \Illuminate\Http\JsonResponse
     */
    public function getRange()
    {
        $adminId = request()->input('admin_id', 110);
        $eId     = request()->input('enterprise_id');
        $isNew   = request()->input('is_new');
        if (empty($adminId)) {
            error(700005);
        }
        if ($isNew) {
            $dataRangeList = (new EnterpriseDataRange())->accountDataRangeByEid($eId);
        } else {
            $dataRangeList = (new EnterpriseDataRange())->accountDataRange($adminId);
        }

        return success($dataRangeList);
    }
}
