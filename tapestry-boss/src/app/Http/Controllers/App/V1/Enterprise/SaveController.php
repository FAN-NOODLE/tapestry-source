<?php

namespace App\Http\Controllers\App\V1\Enterprise;

use App\Common\Constant;
use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Models\Enterprise;
use Illuminate\Http\Request;
use App\Models\EnterpriseApp;
use App\Models\EnterpriseAdmin;
use CrsSdk\Console\Application;
use Illuminate\Support\Facades\DB;
use App\Services\EnterpriseService;
use Illuminate\Support\Facades\Log;
use App\Models\EnterpriseAppInstall;
use App\Http\Requests\EnterpriseRequest;
use App\Services\EnterpriseAdminService;
use App\Http\Requests\EnterpriseAdminRequest;
use App\Http\Controllers\App\V1\BaseController;

class SaveController extends BaseController
{

    /**
     * 企业验证器
     *
     * @var EnterpriseRequest
     */
    protected $enterpriseRequest;

    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * 企业管理员验证器
     *
     * @var EnterpriseAdminRequest
     */
    protected $enterpriseAdminRequest;

    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * SaveController constructor.
     * @param EnterpriseRequest $enterpriseRequest
     * @param EnterpriseService $enterpriseService
     * @param EnterpriseAdminRequest $enterpriseAdminRequest
     * @param EnterpriseAdminService $enterpriseAdminService
     * @throws \Exception
     */
    public function __construct(
        EnterpriseRequest $enterpriseRequest,
        EnterpriseService $enterpriseService,
        EnterpriseAdminRequest $enterpriseAdminRequest,
        EnterpriseAdminService $enterpriseAdminService
    )
    {
        parent::__construct();

        // 企业验证器赋值
        $this->enterpriseRequest = $enterpriseRequest;
        // 企业service赋值
        $this->enterpriseService = $enterpriseService;
        // 企业管理员验证器赋值
        $this->enterpriseAdminRequest = $enterpriseAdminRequest;
        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 编辑企业
     *
     * @param Request $request
     *
     * @return array|mixed
     * @author 贺树鹏 heshupeng@vchangyi.com
     * @date   2019/02/16 38:16
     *
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'enterprise_id',
            'name',
            'contact',
            'mobile',
            'address',
            'crs_status',
            'store_status',
            'contact_status',
            'erp_status',
            'crm_status',
            'corp_id',
            'business_source'
        ]);

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return notice(201101);
        }

        // 切换SHOP域名
        switchAppHost($params['enterprise_id']);

        // 去空
        $params['name']    = trim($params['name']);
        $params['contact'] = trim($params['contact']);
        $params['address'] = trim($params['address']);

        // 请求数据验证
        $validator = $this->enterpriseRequest->common($params);
        if ($validator->fails()) {
            throw new Error(201100, $validator->errors()->first());
        }

        // 企业信息参数处理
        $enterpriseId = $params['enterprise_id'];
        unset($params['enterprise_id']);
        if (empty($params)) {
            return notice(201003);
        }

        // 封装最后操作时间
        $params['edit_at'] = time();

        DB::beginTransaction();
        try {
            // 企业信息保存
            $this->enterpriseService->enterpriseSave($enterpriseId, $params);

            // 企业管理员
            $adminDate             = $request->only(['username', 'account']);
            $adminDate['username'] = trim($adminDate['username']);
            $adminDate['account']  = trim($adminDate['account']);
            if (empty($adminDate['username']) || empty($adminDate['account'])) {
                return notice(202102);
            }

            // 手机号是否已存在
            $mobile = EnterpriseAdmin::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('role', '<>', EnterpriseAdmin::SUPER_ADMIN)
                ->where('account', $adminDate['account'])
                ->first();

            if (!empty($mobile)) {
                return notice(202003);
            }

            // 查找超管ID
            $enterpriseAdmin = EnterpriseAdmin::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('role', EnterpriseAdmin::SUPER_ADMIN)
                ->first();
            // 企业管理员信息保存
            $adminDate['enterprise_id'] = $enterpriseId;
            $this->enterpriseAdminService->enterpriseAdminSave(
                $enterpriseAdmin->id,
                $adminDate
            );

            // 安装应用编辑保存
            $appIds = $request->post('app_ids', []);
            // 首先初始化数据
            EnterpriseAppInstall::query()->where('enterprise_id', $enterpriseId)->delete();

            $enterpriseInfo = Enterprise::query()->find($enterpriseId);

            // 应用ID不为空
            if (!empty($appIds)) {
                // 企业应用标识
                $identity = [];
                foreach ($appIds as $appId) {
                    // 查询应用ID是否存在
                    $appData = EnterpriseApp::query()->where('id', $appId)->first(['id', 'identity']);
                    // 构造应用唯一标识数组，同步信息使用
                    array_push($identity, $appData->identity);

                    if (empty($appData->id ?? '')) {
                        return notice(203003);
                    }

                    // 组装数据
                    $where = [
                        'enterprise_id' => $enterpriseId,
                        'app_id'        => $appId,
                        'identity'      => $appData->identity
                    ];

                    try {
                        // 添加企业应用
                        EnterpriseAppInstall::query()->create($where);
                    } catch (\Exception $e) {
                        return notice(203004);
                    }
                }

                // 开始同步企业应用到 console
                try {
                    $_SERVER['enterprise_hash'] = $enterpriseInfo->is_split_flow == 1 ? $enterpriseInfo->hash : 'public';
                    (new Application())->syncApplication([
                        'enterprise_id' => $enterpriseId,
                        'applications'  => $identity
                    ]);
                } catch (\Exception $e) {
                    Log::info('同步企业应用失败：' . $e->getMessage());
                }
            }


            $hashLockKey = redisKey('ENTERPRISE_HASH_INFO', ['hash' => $enterpriseInfo->hash]);
            $domainLockKey = redisKey('ENTERPRISE_DOMAIN_INFO', ['domain' => $enterpriseInfo->domain]);

            // 企业信息编辑 .删除缓存
            delRedis([
                RedisKey::REDIS_APP_NOVICE . $enterpriseInfo->id,
                RedisKey::REDIS_APP_GROUP . $enterpriseInfo->id,
                RedisKey::REDIS_APP_SECKILL . $enterpriseInfo->id,
                RedisKey::REDIS_ENTERPRISE_INFO . $enterpriseInfo->hash,
                $hashLockKey,
                $domainLockKey
            ]);
            DB::commit();

            return success(['data' => []]);
        } catch (\Exception $e) {
            DB::rollback();
            throw new Error($e->getCode(), $e->getMessage());
        }
    }
}
