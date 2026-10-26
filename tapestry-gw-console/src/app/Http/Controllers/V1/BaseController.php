<?php

namespace App\Http\Controllers\V1;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use App\Library\AccountSdk;
use App\Micros\Authority\AccountMicro;
use App\Models\OperationLog;
use App\Services\RoleService;
use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Shop\StoreWorker;
use Exception;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Route;
use App\Models\Role;
use Illuminate\Support\Facades\Redis;

abstract class BaseController extends Controller
{
    const CLIENT = 'console';

    /**
     * 管理员类型： 超管
     */
    const ADMIN_TYPE_SUPER = 0;

    const GUIDE_MEMBER_ROLE = '导购';

    /**
     * 登录有效期
     */
    const TOKEN_EXPIRE = 86400;

    /**
     * Token解析出来的值
     *
     * @var array|null
     */
    protected $tokenInfo = null;

    /**
     * 企业服务类型
     *
     * @int null
     */
    protected $serveType = 0;

    /**
     * 是否必须登录
     *
     * @var bool
     */
    protected $requireLogin = true;

    /**
     * 当前登录的企业管理员
     *
     * @var null
     */
    protected $enterpriseAdminId = null;

    /**
     * 当前登录的管理员所在的企业ID
     *
     * @var null
     */
    protected $enterpriseId = null;

    /**
     * 当前登录的企业管理员所在的企业信息
     *
     * @var null
     */
    protected $enterprise = null;

    /**
     * 当前登录的门店Id
     *
     * @var null
     */
    protected $storeIds = null;

    /**
     * 当前登录的企业管理员的信息
     *
     * @var array
     */
    protected $enterpriseAdmin = [];

    /**
     * 当前管理员是否为超管0:否,1:是
     *
     * @var int
     */
    protected $isSuperAdmin = 0;

    /**
     * 是否有商城权限 1：是 2：否
     *
     * @var int
     */
    protected $hasShopAuth = 2;

    /**
     * 企业域名
     *
     * @var null|string
     */
    protected $qyDomain = null;

    /**
     * 应用标识
     *
     * @var null|string
     */
    protected $appIdentifier = 'crs-shop';

    /**
     * 应用名称
     */
    private $powerNameBox = '应用';

    /**
     * 应用名称CHILD
     */
    private $powerNameChild = '我的应用';

    public function __construct()
    {
        parent::__construct();

        if (!$this->requireLogin) {
            return true;
        }

        $token = $this->request->header('Auth-Token');

        if ($token === null) {
            throw new Error(450000, '请登录后操作');
        }
        try {
            $tokenInfo = json_decode(Crypt::decrypt($token), true);
        } catch (Exception $exception) {
            throw new Error(450000, 'Auth-Token error: ' . $exception->getMessage());
        }
        //使用redis判断过期时间
        $key = sprintf("gw-scrm:token:%s", $token);
        //15分钟未操作则过期
        if (!Redis::connection()->expire($key, 3600 * 24)) {
            throw new Error(451004, "身份过期", [], 401);
        }

        if (!isset($tokenInfo['is_tmp']) || $tokenInfo['is_tmp'] == 1) {
            throw new Error(451005, '请先换取正式登录token');
        }

        $this->tokenInfo = $tokenInfo;

        // 取出cookie中的企业ID和管理员ID,企业服务类型
        $this->enterpriseId      = $tokenInfo['enterprise_id'];
        $this->enterpriseAdminId = $tokenInfo['admin_id'];
        $this->serveType         = $tokenInfo['serve_type'];
//        $this->storeIds          = $tokenInfo['store_ids'] ?? 0;

        $dataRage = (new AccountMicro(['admin_id' => $this->enterpriseAdminId]))->dataRange();

        $this->storeIds    = $dataRage['store_ids'] ?? [];
        $this->hasShopAuth = $dataRage['has_shop_auth'] ?? 2;

        // 判断BOSS是否关闭应用
        $applicationCloseInfo = Redis::connection()->hgetall('enterprise:application:close:' . $this->enterpriseId);
        if ($applicationCloseInfo) {
            notice(451011, [], $applicationCloseInfo);
        }

        // 调用企业服务接口
        $enterpriseAdmin = (new EnterpriseAdmin())->detail(['admin_id' => $this->enterpriseAdminId]);

        $this->enterpriseAdmin = $enterpriseAdmin['data'];
        $this->enterprise      = $enterpriseAdmin['data']['enterprise'];

        /**
         * 切流量、切db时用到
         */
        $_SERVER['enterprise_hash'] = $this->enterprise['hash'];

        if ($this->enterpriseAdmin['status'] == 2) {
            throw new Error(451000, '您的管理员身份已被禁用');
        }

        $this->qyDomain = $this->enterprise['domain'];

        // 判断当前管理员是否为超管
        if ($this->enterpriseAdmin['is_super'] == 1) {
            $this->isSuperAdmin = 1;
        } else {
            // 判断权限
            try {
                $result = (new AccountSdk())->accountAccessControl(['admin_id' => $this->enterpriseAdminId, 'url' => $this->request->path()]);
                if (!isset($result['access']) || $result['access'] == false) {
                    throw new Error(451006, '无权限访问');
                }
            } catch (Exception $exception) {
                throw new Error(451006, 'Auth-Token error: ' . $exception->getMessage());
            }
        }

        if (is_array($params = request()->input())) {
            Log::error('操作日志', [
                'params' => $params,
                'url'=> request()->url(),
                'account' => $this->enterpriseAdmin
            ]);
        }

        return true;
    }

    /**
     * 获取企业配置
     * @param string $enterpriseConfigMenu
     * @return mixed
     */
    public function enterpriseConfig(string $enterpriseConfigMenu)
    {
        $redisKey = sprintf(Constant::ENTERPRISE_CONFIG, $this->enterpriseId);
        $redis    = Redis::connection();
        $data     = $redis->get($redisKey);
        if (empty($data)) {
            $enterpriseSdk = new EnterpriseConfig();
            $config        = $enterpriseSdk->info(['enterprise_id' => $this->enterpriseId])['data'];
            $data          = array_column($config, 'value', 'key');
            $redis->setex($redisKey, 60, json_encode($data));
        } else {
            $data = json_decode($data);
        }
        return $data[$enterpriseConfigMenu] ?? null;
    }

    /**
     * 写入操作日志
     *
     * @param $content string 操作内容
     *
     * @return bool
     */
    public function writeOperatorLog($content)
    {
        if (!$this->requireLogin) {
            return false;
        }

        $path      = '/' . $this->request->path();
        $routeInfo = $this->searchRoute($path);
        if (!empty($routeInfo)) {
            // 请求参数进行加密
            $ciphertext = [
                'app_key' => config('app.key'),
                'text'    => Crypt::encrypt(json_encode($this->request->all(), JSON_UNESCAPED_UNICODE))
            ];
            OperationLog::query()->create([
                'enterprise_id' => $this->enterpriseId,
                'admin_id'      => $this->enterpriseAdminId,
                'serve_type'    => $this->serveType,
                'admin_name'    => $this->enterpriseAdmin['username'],
                'module'        => $routeInfo['describe'],
                'page'          => $routeInfo['page'],
                'content'       => $content,
                'ip'            => $this->request->getClientIp(),
                'user_agent'    => $this->request->userAgent(),
                'http_method'   => strtoupper($this->request->method()),
                'http_params'   => json_encode($ciphertext, JSON_UNESCAPED_UNICODE)
            ]);
        }

        return true;
    }

    /**
     * 返回路由详情
     *
     * @param $route string 完整路由
     *
     * @return array|null
     */
    private function searchRoute($route)
    {
        $routes = Route::getRoutes();

        // 匹配获取路由详情
        foreach ($routes as $key => $value) {
            if ($value['uri'] === $route) {
                return [
                    'method'   => $value['method'],
                    'route'    => $value['uri'],
                    'page'     => $value['action']['page'] ?? '',
                    'describe' => $value['action']['describe'] ?? ''
                ];
            }
        }

        return null;
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

    protected function getAuthorityDataForV1($enterpriseInfo, $adminId)
    {
        $enterpriseId = $enterpriseInfo['enterprise_id'];
        $serveType    = $enterpriseInfo['serve_type'];
        $role         = $enterpriseInfo['role'];
        $account      = $enterpriseInfo['account'] ?? '';
        $mobile       = $enterpriseInfo['account_mobile'] ?? '';

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
                        'mobile'        => $accountArr[1] ?? $mobile
                    ]);
                }
            } catch (Exception $e) {
                Log::info('admin login get stores info fail :' . json_encode($e->getMessage()));
            }

            if (!empty($storeIds)) {
                // todo::临时给店长超管权限
                $role  = 0;
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
        $data['power']  = $power;
        $data['powers'] = null;
        // 不是超管返回前端加密数据
        if (!$this->isSuperAdmin) {
            $roleInfo       = $this->getStores($role);
            $data['powers'] = $roleInfo['powers'];
        }

        // 服务类型
        $data['serve_type'] = $serveType;

        return $data;
    }

    /** 导出csv
     *
     * 要导出的数组格式的数据
     * @param $dataList
     * 导出的Excel数据第一列表头
     * @param $headList
     * 输出Excel表格文件名
     * @param $fileName
     * 直接输出到浏览器or输出到指定路径文件下，分两种情况。
     * $exportUrl = 'php://output' 表示表示直接输出到浏览器自动下载。
     * $exportUrl = "服务器目录地址/文件名.csv" 表示输出到指定路径文件下。
     * 举例：$exportUrl = "/data/a.csv")
     * @param $exportUrl
     */
    public function toExcel($dataList, $headList, $fileName, $exportUrl = 'php://output')
    {
        header("Content-Description: File Transfer");
        header("Content-type:application/csv;charset=UTF-8");
        header('Content-Disposition: attachment; filename="' . basename($fileName) . '"');
        header('Expires: 0');//缓存过期时间
        header('Cache-Control: must-revalidate, post-check=0, pre-check=0');
        header('Pragma: public');

        $fp = fopen($exportUrl, 'a+');
        //输出Excel列名信息
        foreach ($headList as $key => $value) {
            //CSV的Excel支持GBK编码，一定要转换，否则乱码
//            $headList[$key] = iconv('utf-8', 'gbk', $value);
            $headList[$key] = mb_convert_encoding($value, "GBK", "UTF-8");
        }
        //将数据通过fputcsv写到文件句柄
        if (!empty($headList)) {
            fputcsv($fp, $headList);
        }

        //计数器
        $num = 0;

        //每隔$limit行，刷新一下输出buffer，不要太大，也不要太小
        $limit = 100000;

        //逐行取出数据，不浪费内存
        $count = count($dataList);
        for ($i = 0; $i < $count; $i++) {

            $num++;

//        刷新一下输出buffer，防止由于数据过多造成问题
            if ($limit == $num) {
                ob_flush();
                flush();
                $num = 0;
            }

            $row = $dataList[$i];
            foreach ($row as $key => $value) {
//                $row[$key] = iconv('utf-8', 'gbk', $value);
                $row[$key] = mb_convert_encoding($value, "GBK", "UTF-8");
            }
            fputcsv($fp, $row);
        }
        fclose($fp);

        if ($exportUrl != 'php://output') {
            ob_clean();
            flush();
            readfile($exportUrl);
        }
    }
}
