<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波、赵杰
 * Date: 2020/4/10
 * Time: 17:06
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\MicroScrm\GoMicroScrm;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Cache;
use Carbon\Carbon;

/**
 * 店员详情
 * Class DetailController
 *
 * @package App\Http\Controllers\App\V1\Member
 */
class DetailController extends Controller
{
    /**
     * 获取店员详情
     *
     * @param Request     $request
     * @param StoreWorker $storeWorker
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Exception
     */
    public function __invoke(Request $request)
    {
        // 参数校验
        $rules   = [
            'store_id' => 'integer',
        ];
        $message = [
            'store_id.integer' => '门店ID参数类型有误',
        ];
        $this->validate($request, $rules, $message);
        $storeId       = $request->get('store_id', 0);
        $id            = $request->get('id', 0);
        $loginCacheKey = $request->get('login_cache_key', '');

        $wheres      = [
            'enterprise_id' => $this->enterpriseId,
            'store_id'      => $storeId,
            'mobile'        => $this->mobile,
            'worker_status' => Constant::WORKER_ONLINE,
            'worker_id'     => $id
        ];
        $storeWorker = new StoreWorker();
        try {
//        $detail = $storeWorker->detail($wheres);

            $detail = GoMicroScrm::getMemberInfo([
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->mobile,
            ]);
//        dd($detail, $memberInfo);
        } catch (\Exception $e) {
            Log::error("调用接口Shop:StoreWork->detail出错，错误信息：{$e->getMessage()}, 错误code：{$e->getCode()}", $wheres);
            error($e->getCode());
        }
        if (empty($detail)) {
            notice('store_no_authority');
        }
        $storeAuthority = [];
        $guideRole      = $this->role;
        $guideRoleNew = Constant::LOGIN_ROLE_STORE_MEMBER;
        if (Constant::LOGIN_MANAGER_ROLE_WORK_GUIDE == $detail['type']) {
            $guideRoleNew = Constant::LOGIN_ROLE_STORE_MANAGER;
        }

        if ($storeId) {
            // 记录导购最新登录门店
            $httpParam['where'] = [
                'mobile'        => $this->mobile,
                'enterprise_id' => $this->enterpriseId
            ];
            $httpParam['data']  = [
                'recently_logged_store_id' => $storeId,
                'recently_logged_time'     => time()
            ];
            $storeWorker->update($httpParam);
        }

        // 切换门店后店员/店长角色切换
        if ($guideRole != $guideRoleNew) {
            $guideRole = $guideRoleNew;

            // 兼容前端老代码,防止报错
            if ($loginCacheKey) {
                $cacheData         = Cache::get($loginCacheKey);
                $cacheData['role'] = $guideRoleNew;

                $expiredTime = Carbon::createFromTimestamp($cacheData['expiredTime']);
                Cache::forget($loginCacheKey);
                Cache::add($loginCacheKey, $cacheData, $expiredTime);
            }
        }

        $result = [
            'id'                => $detail['id'],
            'enterprise_id'     => $this->enterpriseId,
            'name'              => $detail['name'],
            'mobile'            => $detail['mobile'],
            'worker_no'         => $detail['userid'],
            'type'              => $detail['type'],
            'qr_code'           => $detail['qr_poster_code'],
            'worker_status'     => $detail['status'] == 1 ? 1 : 2, // 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业
            'avatar'            => $detail['avatar'],
            'guide_permissions' => $storeAuthority,
            'guide_role'        => $guideRole,
        ];
        return success($result);
    }
}
