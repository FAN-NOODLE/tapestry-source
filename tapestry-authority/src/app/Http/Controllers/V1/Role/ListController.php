<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers\V1\Role;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseRoleService;


class ListController extends BaseController
{
    public function index(
        EnterpriseRoleService $enterpriseRoleService,
        EnterpriseAdminService $enterpriseAdminService
    )
    {
        $params = $this->request->input();

        if (empty($params['enterprise_id'])) {
            notice(700014);
        }

        $page         = $params['page'] ?? Constant::DEFAULT_PAGE;
        $size         = $params['size'] ?? Constant::DEFAULT_LIMIT;
        $sortType     = $params['sort_type'] ?? '';
        $enterpriseId = $params['enterprise_id'];

        $extParams = [];
        if (isset($params['admin_name']) && !empty($params['admin_name'])) {
            // 根据管理员名称和企业ID，搜索管理员id
            $adminWhere = [
                ['enterprise_id', '=', $enterpriseId],
                ['username', 'LIKE', '%' . $params['admin_name'] . '%']
            ];
            $adminList  = $enterpriseAdminService->getAdminListByWhere($adminWhere);
            if (!empty($adminList)) {
                $adminList                = $adminList->toArray();
                $extParams['operator_id'] = array_column($adminList, 'id');
            }

            if (isset($extParams['operator_id']) && empty($extParams['operator_id'])) {
                return success([
                    'page'  => $page,
                    'size'  => $size,
                    'total' => 0,
                    'list'  => []
                ]);
            }
        }

        $data  = [];
        $where = [['enterprise_id', '=', $enterpriseId]];
        // 搜索角色
        if (isset($params['role_name']) && !empty($params['role_name'])) {
            $where[] = ['role_name', 'LIKE', '%' . $params['role_name'] . '%'];
        }

        // 排序
        $extParams['sort_type'] = $sortType;
        $result = $enterpriseRoleService->getRoleListByWhere($where, $page, $size, $extParams);
        if (!empty($result)) {
            $data = [
                'page'  => $page,
                'size'  => $size,
                'total' => $result->total(),
                'list'  => $result->items()
            ];
        }

        return success($data);
    }
}
