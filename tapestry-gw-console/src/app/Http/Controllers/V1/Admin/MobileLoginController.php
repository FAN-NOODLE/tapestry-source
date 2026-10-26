<?php

namespace App\Http\Controllers\V1\Admin;

use CrsSdk\Work\StoreWorker;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Crypt;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Models\OperationLog;
use App\Services\RoleService;
use CrsSdk\Boss\EnterpriseAdmin;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class MobileLoginController extends BaseController
{
    /**
     * 管理员类型： 超管
     */
    const ADMIN_TYPE_SUPER = 0;
    /**
     *  服务状态： 开启
     */
    const SERVE_OPEN_TRUE = 1;
    /**
     *  服务状态： 禁用
     */
    const SERVE_OPEN_FALSE = 2;
    /**
     *  服务类型： 商城
     */
    const SERVE_TYPE_SHOP = 1;
    /**
     *  服务类型： crs
     */
    const SERVE_TYPE_CRS = 2;
    /**
     *  服务类型： crs+商城
     */
    const SERVE_TYPE_ALL = 0;

    /**
     * 服务类型：crs+商城都没有
     */
    const SERVE_TYPE_NONE = 3;

    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业ID
     */
    protected $enterpriseId = null;

    /**
     * 应用名称
     */
    private $powerNameBox = '应用';

    /**
     * 应用名称CHILD
     */
    private $powerNameChild = '我的应用';

    /**
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        return success(['废弃']);

        $mobile = trim($this->request->post('mobile'));
        $code   = trim($this->request->post('code'));

        if (empty($mobile)) {
            // 手机号不能为空
            error(451013);
        }

        if (empty($code)) {
            // 验证码不能为空
            error(451016);
        }

        // 缓存中获取验证码
        $accountArr = explode('-', $mobile);
        $cacheCode  = Redis::connection()->get('easy-sms:mobile:' . $accountArr[1]);
        if (empty($cacheCode)) {
            // 验证码已失效
            error(451017);
        }

        if ($cacheCode != $code) {
            // 验证码不正确
            error(451018);
        }

        // 调用服务接口，查询企业管理员是否存在，以及账号密码的有效性
        $result             = $this->checkAdmin($mobile);
        $this->enterpriseId = $result['enterprise_id'];

        // 删除关闭应用标识缓存
        Redis::connection()->del(['enterprise:application:close:' . $this->enterpriseId]);

        $adminId = $result['id'];

        // 获取服务类型
        $serveType   = $this->getEnterpriseServeType($result);
        $RoleService = new RoleService($this->enterpriseId);

        // 获取当前管理员权限
        if ($result['role'] == self::ADMIN_TYPE_SUPER) {
            $power = $RoleService->getAllPower();
        } else {
            $power = $RoleService->getPower($result['role']);
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
        $result['power']  = $power;
        $result['powers'] = null;
        // 不是超管返回前端加密数据
        if (!$this->isSuperAdmin) {
            $roleInfo         = $this->getStores($result['role']);
            $result['powers'] = $roleInfo['powers'] ?? null;
        }

        // 服务类型
        $result['serve_type'] = $serveType;

        // 通过手机号获取管理的门店ID列表
        $storeIds = (new StoreWorker())->getStoreIds([
            'enterprise_id' => $this->enterpriseId,
            'mobile'        => $accountArr[1]
        ]);

        Log::info('店长管理员管理的门店列表：' . json_encode($storeIds));

        $token = Crypt::encrypt(json_encode([
            'admin_id'      => $adminId,
            'enterprise_id' => $this->enterpriseId,
            'serve_type'    => $serveType,
            'store_ids'     => $storeIds ?? [],
            'expire_time'   => time() + parent::TOKEN_EXPIRE
        ]));

        // 写入日志(登陆接口使用不了BaseController写入日志方法,需单独处理)
        OperationLog::query()->create([
            'enterprise_id' => $this->enterpriseId,
            'admin_id'      => $adminId,
            'serve_type'    => $serveType,
            'admin_name'    => $result['username'],
            'module'        => '管理员管理',
            'page'          => '系统登陆页',
            'content'       => '系统登陆',
            'ip'            => $this->request->getClientIp(),
            'user_agent'    => $this->request->userAgent(),
            'http_method'   => strtoupper($this->request->method()),
            'http_params'   => json_encode($this->request->all(), JSON_UNESCAPED_UNICODE)
        ]);

        // 登陆成功删除缓存验证码
        Redis::connection()->del(['easy-sms:mobile:' . $accountArr[1]]);

        return success([
            'data'  => $result,
            'token' => $token,
            'hash'  => $result['enterprise']['is_split_flow'] ? $result['enterprise']['hash'] : 'public'
        ]);
    }

    /**
     * 调用服务接口，检查管理员信息是否有效
     *
     * @param $account
     *
     * @return mixed
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/28 11:30 AM
     */
    private function checkAdmin($account)
    {
        $entAdmin = new EnterpriseAdmin();

        $response = $entAdmin->detail(['account' => $account]);
        $data     = $response['data'];
        // 管理员被禁用
        if ($data['status'] == 2) {
            throw new Error(10000, '您的管理员身份已被禁用');
        }

        if ($data['account'] === $account) {
            // 更新最后登录时间
            $entAdmin->loginUpdate(['admin_id' => $data['id']]);
            return $data;
        }

        throw new Error(451003, '管理员账号或者密码错误');
    }

    /**
     * 获取当前企业开通服务类型
     *
     * @param $enterpriseInfo
     *
     * @return int
     */
    public function getEnterpriseServeType($enterpriseInfo)
    {
        // 初始化企业服务了开启类型
        $serveType = self::SERVE_TYPE_ALL;

        // 仅开通商城
        if (($enterpriseInfo['enterprise']['store_status'] == self::SERVE_OPEN_TRUE) && ($enterpriseInfo['enterprise']['crs_status'] == self::SERVE_OPEN_FALSE)) {
            $serveType = self::SERVE_TYPE_SHOP;
        }

        // 仅开通CRS
        if (($enterpriseInfo['enterprise']['store_status'] == self::SERVE_OPEN_FALSE) && ($enterpriseInfo['enterprise']['crs_status'] == self::SERVE_OPEN_TRUE)) {
            $serveType = self::SERVE_TYPE_CRS;
        }

        // 开通crs和商城
        if (($enterpriseInfo['enterprise']['store_status'] == self::SERVE_OPEN_TRUE) && ($enterpriseInfo['enterprise']['crs_status'] == self::SERVE_OPEN_TRUE)) {
            $serveType = self::SERVE_TYPE_ALL;
        }

        // 什么都不开通
        if (($enterpriseInfo['enterprise']['store_status'] == self::SERVE_OPEN_FALSE) && ($enterpriseInfo['enterprise']['crs_status'] == self::SERVE_OPEN_FALSE)) {
            $serveType = self::SERVE_TYPE_NONE;
        }

        return $serveType;
    }

}
