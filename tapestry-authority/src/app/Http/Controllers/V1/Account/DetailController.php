<?php

namespace App\Http\Controllers\V1\Account;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Modules\Rbac\AdminRole;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use App\Services\EnterpriseRoleService;
use CrsSdk\Shop\Store;
use CrsSdk\Shop\StoreArea;
use CrsSdk\Work\Department;
use Illuminate\Support\Facades\Log;

/**
 * 账号详情
 * Class DetailController
 * @package App\Http\Controllers\V1\Account
 */
class DetailController extends BaseController
{
    /**
     * 账号详情
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @api {domain}/micro-authority/v1/account/detail
     * @author lijian lijian@vchangyi.com 2020/04/16
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'mobile',
            'enterprise_id',
            'qy_domain',
            'only'
        ]);

        $input['only'] = $input['only'] ?? 1;

        // 验证
        AccountRequest::detail($input);

        $adminDetail = [];
        // 获取账号详情-及关联的范围和角色
        if (isset($input['id'])) {
            $adminDetail = (new EnterpriseAdminService())->getEnterpriseAdminById($input) ?? [];
        }

        if (isset($input['mobile'])) {
            $adminDetail = (new EnterpriseAdminService())->getEnterpriseAdminByMobile($input) ?? [];
        }

        if (empty($adminDetail)) {
            notice(701006);
        }

        $input['id'] = $adminDetail['id'];

        // 此参数控制返回简单详情和复杂详情
        if ($input['only'] != 1) {
            return success($adminDetail);
        }

        // 获取管理员账号对应的角色
        $adminRole = (new EnterpriseAdminRoleService())->getAdminRoleByAdminId($input['id']);
        $roleIds   = array_column($adminRole, 'role_id');
        if (!empty($roleIds)) {
            $roleList     = (new EnterpriseRoleService())->getRoleList($input['enterprise_id'], $roleIds);
            $roleListData = [];
            array_walk($roleList, function ($value) use (&$roleListData) {
                $roleListData[] = [
                    'id'        => $value['id'],
                    'role_name' => $value['role_name']
                ];
            });

        }

        if ($adminDetail['range_type'] == EnterpriseAdmin::RANGE_TYPE_ONLY) {
            // 获取管理员账号对应的数据范围
            $adminDataRange = (new EnterpriseDataRangeService())->getDataRangeListByAdminId($input['id']);

            $range    = [];
            $storeIds = [];
            $areaIds  = [];
            $depIds   = [];
            foreach ($adminDataRange as $value) {
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

            // 获取区域列表
            if (!empty($areaIds)) {
                $storeAreaParam = [
                    'enterprise_id' => $input['enterprise_id'],
                    'area_ids'      => $areaIds
                ];
                $areaList       = (new StoreArea())->areaListByArea($storeAreaParam);
                $areaList       = array_column($areaList['list'], null, 'id');
            }

            // 获取门店列表
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
                    $list      = (new Department())->departmentList($condition);
                    $depList   = array_column($list['list'], 'dpName', 'dpId');
                }
            } catch (\Exception $exception) {
                Log::info('获取部门列表失败' . $exception->getMessage());
            }

            foreach ($adminDataRange as $rangeV) {
                $pid  = 0;
                $name = '';
                switch ($rangeV['type']) {
                    case EnterpriseDataRange::TYPE_ALL:
                        $name = "全部门店";
                        break;
                    case EnterpriseDataRange::TYPE_SHOP:
                        $name = "小程序";
                        break;
                    case EnterpriseDataRange::TYPE_AREA:
                        $name = isset($areaList[$rangeV['object_id']]) ? $areaList[$rangeV['object_id']]['name'] : '';
                        $pid  = isset($areaList[$rangeV['object_id']]) ? $areaList[$rangeV['object_id']]['pid'] : 0;
                        break;
                    case EnterpriseDataRange::TYPE_STORE:
                        $name = isset($storeList[$rangeV['object_id']]) ? $storeList[$rangeV['object_id']]['name'] : '';
                        $pid  = isset($storeList[$rangeV['object_id']]) ? $storeList[$rangeV['object_id']]['area_id'] : 0;
                        break;
                    case EnterpriseDataRange::TYPE_DEPARTMENT:
                        $name = isset($depList[$rangeV['object_id']]) ? $depList[$rangeV['object_id']] : '';
                        $pid  = 0;
                        break;
                    default:
                        $item['name'] = '';
                        break;
                }
                $range[] = [
                    'id'   => $rangeV['object_id'],
                    'type' => $rangeV['type'],
                    'name' => $name,
                    'pid'  => $pid
                ];
            }
        } else {
//            $range[] = [
//                'id' => 0,
//                'type' => -1,
//                'name' => '不限',
//                'pid'  => 0
//            ];
            $range = [];
        }

        $adminDetail['role_list'] = $roleListData ?? [];
        $adminDetail['range']     = $range;
        $backData                 = array_only($adminDetail, [
            'id',
            'enterprise_id',
            'username',
            'account',
            'email',
            'mobile',
            'password',
            'status',
            'created_at',
            'updated_at',
            'admin_available_system',
            'range_type',
            'role_list',
            'range'
        ]);

        $roleList = (new AdminRole($input['id']))->roles();

        $backData['new_role_list']          = $roleList;
        $backData['admin_available_system'] = json_decode($backData['admin_available_system'], true);

        return success($backData);
    }

}