<?php


namespace App\Http\Controllers\V1\Admin;


use App\Http\Controllers\V1\BaseController;
use App\Models\Role;
use App\Services\RoleService;
use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AutoAccountController extends BaseController
{
    public $requireLogin = false;
    const CLIENT = 'console';
    /**
     * 应用名称
     */
    private $powerNameBox = '应用';
    /**
     * 应用名称CHILD
     */
    private $powerNameChild = '我的应用';
    /**
     * 管理员类型： 超管
     */
    const ADMIN_TYPE_SUPER = 0;
    /**
     * 当前管理员是否为超管0:否,1:是
     *
     * @var int
     */
    protected $isSuperAdmin = 0;

    public function getInfoV1(Request $request)
    {
        $adminId = trim($request->admin_id);
        $enterpriseId = trim($request->enterprise_id);
        $serveType = trim($request->serve_type);
        $role = trim($request->role);
        $account = empty(trim($request->account)) ?? "";
        $mobile = empty(trim($request->account_mobile)) ?? "";
        $RoleService = new RoleService($enterpriseId);
        // 通过手机号获取管理的门店ID列表
        $accountArr = explode('-', $account);

        $response = (new EnterpriseAdmin())->detail(['admin_id' => $adminId]);

        $data = $response['data'];

        // 获取当前管理员权限
        if ($role == self::ADMIN_TYPE_SUPER) {
            $power = $RoleService->getAllPower();
            /**
             * 查询店长角色ID todo::临时方案，权限功能完善后，再做修改
             */
        } else {
            try {
                // 返回门店ID列表
                if ($role != self::ADMIN_TYPE_SUPER) {
                    $storeIds = (new StoreWorker())->getStoreIds([
                        'enterprise_id' => $enterpriseId,
                        'mobile' => $accountArr[1] ?? $mobile
                    ]);
                }
            } catch (\Exception $e) {
                Log::info('admin login get stores info fail :' . json_encode($e->getMessage()));
            }

            if (!empty($storeIds)) {
                // todo::临时给店长超管权限
                $role = 0;
                $power = $RoleService->getAllPower();
            } else {
                $power = $RoleService->getPower($role);
            }
        }

        if ($power['shop']) {
            foreach ($power['shop'] as $key => $value) {
                if ($value['name'] == $this->powerNameBox) {
                    foreach ($value['child'] as $k => $v) {
                        if ($v['name'] == $this->powerNameChild) {
                            $power['shop_power'] = $v['child'];
                        }
                    }
                }
            }
        }
        $data['power'] = $power;
        $data['powers'] = null;
        // 不是超管返回前端加密数据
        if (!$this->isSuperAdmin) {
            $roleInfo = $this->getStores($role);
            $data['powers'] = $roleInfo['powers'];
        }

        // 服务类型
        $data['serve_type'] = $serveType;

        return success(["data" => $data]);
    }
    /**
     * 获取角色及关联门店
     *
     * @param $roleId int 角色ID
     * @param $power bool 是否需要前端powers
     *
     * @return array
     */
    public function getStores($roleId, $power = true)
    {
        $data = [
            'stores'         => [],
            'type'           => 1,
            'is_super_admin' => $this->isSuperAdmin,
            'powers'         => null
        ];
        Log::debug('getStores', $data);
        if (!$data['is_super_admin']) {
            $role = Role::query()->select([
                'id',
                'name',
                'describe',
                'type',
                'powers'
            ])->where('id', '=', $roleId)->first();

            // 角色存在
            if (!empty($role)) {
                // 查询角色下的门店
                $role->stores;
                $role           = $role->toArray();
                $role['stores'] = array_column($role['stores'], 'store_id');

                $data = array_merge($data, $role);

                Log::debug('合并array', $data);
            }
        }
        // 前端权限参数，不需要就false
        if (!$power) {
            // unset($data['powers']);
        }

        return $data;
    }
}
