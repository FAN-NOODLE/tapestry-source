<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\EnterpriseAdmin;
use App\Modules\Rbac\Role;
use Illuminate\Http\Request;
use App\Modules\Rbac\AdminRole;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class AccountAccessController extends BaseController
{
    public function index(Request $request)
    {
        $adminId = trim($request->post('admin_id', 0));
        $url     = trim($request->post('url'));
        Log::info('request url:' . $url);

        $access  = false;
        // 判断是否是超管
        $adminInfo = EnterpriseAdmin::query()->where(['id' => $adminId])->first()->toArray();
        if (!empty($adminInfo) && $adminInfo['is_super'] == 1) {
            return success(['access' => true]);
        }
        // 非超管，判断是否有访问权限
        $apiJson = $this->redis->get($adminId);
        $apiList = empty($apiJson) ? [] : json_decode($apiJson, true);
        if (in_array($url, $apiList)) {
            $access = true;
        }

        return success(['access' => $access]);
    }
}
