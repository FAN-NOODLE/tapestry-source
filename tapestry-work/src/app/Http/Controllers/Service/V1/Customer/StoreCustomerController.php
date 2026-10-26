<?php
/**
 * 外部联系人列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/3/24
 * Time: 16:46
 */
namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class StoreCustomerController extends BaseController
{
    public function index(Request $request)
    {
        $store = $request->post('store', []);
        $enterpriseId = $request->post('enterprise_id');

        if (!$enterpriseId) {
            error(600000);
        }
        $userList = $this->getData($enterpriseId, $store);
        if (empty($userList)) {
            $ecIds = [];
        } else {
            $ecIds = array_unique(array_column($userList, 'ec_id'));
        }

        return success(['data' => $ecIds]);
    }

    /**
     * 获取数据权限内的外部联系人
     * @param $enterpriseId
     * @param $storeRole
     * @return mixed
     */
    public function getData($enterpriseId, $storeRole)
    {
        $userList = MemberUserRelation::query()
                ->when(isset($storeRole['is_super_admin'])
                    && ($storeRole['is_super_admin'] != 1), function ($query) use ($storeRole) {
                        $query->when(isset($storeRole['type'])
                            && $storeRole['type'] == self::DATA_PERMISSION_TYPE_NATURE, function ($query) {
                                $query->{'where'}('store_id', '=', '');
                            })
                        ->when(isset($storeRole['type'])
                            && $storeRole['type'] == self::DATA_PERMISSION_TYPE_STORE
                            && !empty($storeRole['stores']), function ($query) use ($storeRole) {
                                $query->{'whereIn'}('store_id', $storeRole['stores']);
                            })
                        ->when(isset($storeRole['type'])
                            && $storeRole['type'] == self::DATA_PERMISSION_TYPE_ALL
                            && !empty($storeRole['stores']), function ($query) use ($storeRole) {
                                $query->{'whereIn'}('store_id', $storeRole['stores'])
                                    ->{'orWhere'}('store_id', '=', '');
                            });
                    })
            ->where('enterprise_id', '=', $enterpriseId)
            ->limit(15000) // 最多只允许查15000条，否则MySQL会查不出来
            ->get(['ec_id'])
            ->toArray();
        return $userList;
    }
}
