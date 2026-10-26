<?php


namespace App\Http\Controllers\V1\Account;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\AdminRoleRelation;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use CrsSdk\Shop\Store;
use CrsSdk\Shop\StoreArea;
use CrsSdk\Work\Department;
use Illuminate\Support\Facades\Log;

/**
 * 账号列表
 * Class ListController
 * @package App\Http\Controllers\V1\Account
 */
class ListController extends BaseController
{
    /**
     * 搜索类型：username-账号名搜索 mobile-电话号码搜索
     */

    /**
     * 账号列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @api {domain}/micro-authority/v1/account/list
     * @author lijian lijian@vchangyi.com 2020/04/16
     */
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'search_type',
            'keyword',
            'object_id',
            'role_id',
            'page',
            'size',
            'qy_domain'
        ]);

        // 验证
        AccountRequest::list($input);

        $search['page']          = $input['page'] ?? Constant::DEFAULT_PAGE;
        $search['size']          = $input['size'] ?? Constant::DEFAULT_LIMIT;
        $search['enterprise_id'] = $input['enterprise_id'];

        // 只传入门店ID 需要调用门店获取该门店所有层级得ID
        $rangeAdmins = [];
        $areaIds     = [];
        if (isset($input['object_id']) && !empty($input['object_id'])) {
            // 商城不用查
            if ($input['object_id'] != $input['enterprise_id']) {
                // 使用门店ID获取该门店所属所有层级的区域ID
                $storeParams = [
                    'enterprise_id' => $input['enterprise_id'],
                    'store_ids'     => [$input['object_id']],
                ];
                $areaSet     = (new StoreArea())->areaListByStore($storeParams);
                $areaIds     = !empty($areaSet) ? $areaSet[0]['area_ids'] : [];
            }
            // 合并传参门店ID到区域ID集合中
            $areaIds[]          = $input['object_id'];
            $input['object_id'] = array_values(array_unique($areaIds));
            // 获取区域门店集合关联的账号ID
            $dataRangeList = (new EnterpriseDataRangeService())->getDataRangeList($input['enterprise_id'], $input['object_id']);
            $rangeAdmins   = array_column($dataRangeList, 'enterprise_admin_id');
            if (empty($rangeAdmins)) {
                return success([
                    'list'  => [],
                    'page'  => $search['page'],
                    'size'  => $search['size'],
                    'total' => 0

                ]);
            }
            $search['ids'] = $rangeAdmins;
        }

        // 如果有传入角色ID 则先获取角色表管理员ID集合
        $roleAdminS = [];
        if (isset($input['role_id']) && !empty($input['role_id'])) {
            $roleAdminS = AdminRoleRelation::adminRelationIds($input['role_id']);
            if (empty($roleAdminS)) {
                return success([
                    'list'  => [],
                    'page'  => $search['page'],
                    'size'  => $search['size'],
                    'total' => 0

                ]);
            }
            $search['ids'] = $roleAdminS;
        }

        // 如果二者同时存在 则取交集
        if (isset($input['role_id']) && !empty($input['role_id']) && isset($input['object_id']) && !empty($input['object_id'])) {
            $search['ids'] = array_intersect($rangeAdmins, $roleAdminS);
            if (empty($search['ids'])) {
                return success([
                    'list'  => [],
                    'page'  => $search['page'],
                    'size'  => $search['size'],
                    'total' => 0

                ]);
            }
        }

        // 搜索主表 账号表
        if ($input['search_type'] == EnterpriseAdminService::SEARCH_MOBILE) {
            $search['mobile'] = $input['keyword'];
        }

        if ($input['search_type'] == EnterpriseAdminService::SEARCH_ACCOUNT) {
            $search['username'] = $input['keyword'];
        }

        // 获取新版的角色绑定的管理员, 覆盖上面的
        if (!empty($input['role_id'])) {
            $search['ids'] = AdminRoleRelation::adminRelationIds($input['role_id']);
        }

        // 获取管理员账号列表
        $list     = (new EnterpriseAdminService())->listEnterpriseAdmin($search)->toArray();
        $adminIds = array_column($list['data'], 'id');
        if (empty($adminIds)) {
            return success([
                'list'  => [],
                'page'  => $search['page'],
                'size'  => $search['size'],
                'total' => 0

            ]);
        }

        // 获取数据范围
        $storeIds      = [];
        $areaIds       = [];
        $depIds        = [];
        $dataRangeList = (new EnterpriseDataRangeService())->getDataRangeList($input['enterprise_id'], null, $adminIds);
        foreach ($dataRangeList as $value) {
            if ($value['type'] == EnterpriseDataRange::TYPE_AREA) {
                $areaIds[] = $value['object_id'];
            }
            if ($value['type'] == EnterpriseDataRange::TYPE_STORE) {
                $storeIds[] = $value['object_id'];
            }

            if ($value['type'] == EnterpriseDataRange::TYPE_DEPARTMENT) {
                $depIds[] = $value['object_id'];
            }
        }

        // 根据区域ID获取区域列表
        $areaList = [];
        if (!empty($areaIds)) {
            $storeAreaParam = [
                'enterprise_id' => $input['enterprise_id'],
                'area_ids'      => $areaIds
            ];
            $areaList       = (new StoreArea())->areaListByArea($storeAreaParam);
            $areaList       = array_column($areaList['list'], null, 'id');
        }

        // 根据门店ID获取门店列表
        $storeList = [];
        if (!empty($storeIds)) {
            $storeParam = [
                'enterprise_id' => $input['enterprise_id'],
                'store_ids'     => $storeIds,
                'page'          => Constant::DEFAULT_PAGE,
                'size'          => 999999
            ];
            $storeList  = (new Store())->lists($storeParam);
            $storeList  = array_column($storeList['list'], null, 'id');
        }

        // 获取部门列表
        $depList = [];
        try {
            if (!empty($depIds) && !empty($input['qy_domain'])) {
                // 获取部门列表
                $condition = [
                    'qy_domain'      => $input['qy_domain'],
                    'app_identifier' => 'contact',
                    'params'         => [
                        'pageSize' => count($depIds),
                        'dpIds'    => $depIds
                    ],
                ];
                $departList      = (new Department())->departmentList($condition);
                $depList = array_column($departList['list'], 'dpName', 'dpId');
            }
        } catch (\Exception $exception) {
            Log::info('获取部门列表失败' . $exception->getMessage());
        }

        // 拼装数据
        foreach ($dataRangeList as $key => &$item) {
            switch ($item['type']) {
                case EnterpriseDataRange::TYPE_ALL:
                    $item['name'] = "全部门店";
                    break;
                case EnterpriseDataRange::TYPE_SHOP:
                    $item['name'] = "小程序";
                    break;
                case EnterpriseDataRange::TYPE_AREA:
                    $item['name'] = isset($areaList[$item['object_id']]) ? $areaList[$item['object_id']]['name'] : '';
                    break;
                case EnterpriseDataRange::TYPE_STORE:
                    $item['name'] = isset($storeList[$item['object_id']]) ? $storeList[$item['object_id']]['name'] : '';
                    break;
                case EnterpriseDataRange::TYPE_DEPARTMENT:
                    $item['name'] = isset($depList[$item['object_id']]) ? $depList[$item['object_id']] : '';
                    break;
                default:
                    $item['name'] = '';
                    break;
            }
        }
        unset($item);


        $backData = [];
        foreach ($list['data'] as $key => $admin) {
            $backData[$key] = [
                'id'                  => $admin['id'],
                'username'            => $admin['username'],
                'account'             => $admin['account'],
                'mobile'              => $admin['mobile'],
                'status'              => $admin['status'],
                'created_at'          => $admin['created_at'],
                'updated_at'          => $admin['updated_at'],
                'range_type'          => $admin['range_type'],
                'range'               => [],
                'roles'               => [],
                'admin_role_relation' => $admin['admin_role_relation']
            ];

            if ($admin['range_type'] == EnterpriseAdmin::RANGE_TYPE_ONLY) {
                // 拼接范围数据
                foreach ($dataRangeList as $tangeKey => $range) {
                    if ($range['enterprise_admin_id'] == $admin['id']) {
                        $backData[$key]['range'][] = [
                            'id'   => $range['object_id'],
                            'type' => $range['type'],
                            'name' => $range['name']
                        ];
                        // 减少下次循环次数
                        unset($dataRangeList[$tangeKey]);
                    }
                }
            } else {
                $backData[$key]['range'][] = [
                    'id'   => 0,
                    'type' => -1,
                    'name' => '不限'
                ];
            }
        }

        return success([
            'list'  => $backData,
            'page'  => $search['page'],
            'size'  => $search['size'],
            'total' => $list['total'] ?? 0
        ]);
    }
}