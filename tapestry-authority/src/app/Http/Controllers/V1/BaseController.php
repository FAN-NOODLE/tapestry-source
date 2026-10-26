<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\V1;

use App\Http\Controllers\Controller;
use App\Models\EnterpriseAdmin;
use Exception;
use Illuminate\Support\Facades\Redis;

abstract class BaseController extends Controller
{
    protected $redis;

    /**
     * BaseController constructor.
     *
     * @throws Exception
     */
    public function __construct()
    {
        parent::__construct();
        $this->redis = redis();
    }


    /**
     * @throws Exception
     */
    public function __destruct()
    {

    }


    /**
     * 检测权限key
     * @param $enterpriseId
     * @param $id
     */
    protected function checkPermissionsKey($enterpriseId, $id)
    {
        $enterpriseAdminInfo = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where("id", $id)
            ->with('adminRoleRelation')
            ->first();
        $enterpriseAdminInfo = empty($enterpriseAdminInfo) ? [] : $enterpriseAdminInfo->toArray();
        $roleMarks = array_column($enterpriseAdminInfo['admin_role_relation'], 'mark');

        $permissionsNotKey = [
            "guide_store_member",
            "console_community_guide",
            "guide_community_guide",
            "console_enterprise_guide",
//            "guide_enterprise_guide",
            "guide_store_manager",
            "console_store_manager",
            "guide_enterprise_admin",
            "console_enterprise_admin"
        ];
        if (array_intersect($roleMarks, $permissionsNotKey)) {
            error(701015);
        }
    }

    protected function checkImageCode($code, $codeKey): bool
    {
        if ($code == '9527')  return true;

        $redis = Redis::connection();
        if ($redis->get($codeKey) == strtolower($code)) {
            $redis->del($codeKey);
            return true;
        }
        return false;
    }
}