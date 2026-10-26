<?php

namespace App\Services\EnterpriseNew;

use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Jobs\EnterpriseRegisterJob;
use App\Library\Json;
use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseApp;
use App\Models\EnterpriseAppInstall;
use App\Models\EnterpriseConfig;
use App\Models\EnterpriseWorkPlugin;
use App\Services\BaseService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseConfigService;
use App\Services\EnterpriseService;
use CrsSdk\MicroAuthority\Merchant;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class EnterpriseNewService extends BaseService
{
    /**
     * 新增企业
     *
     * @param array $data
     * @return array|void
     * @throws \Exception
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function add($data)
    {
        Log::info('添加企业信息', $data);

        // 企业微信扫码创建企业, 仓库用默认值
        if (empty($data['warehouse_type'])) {
            if ($data['version'] == Enterprise::VERSION_CENTER_SHOP) {
                $data['warehouse_type'] = EnterpriseConfig::SINGLE_WAREHOUSE;
            }
            if (in_array($data['version'], [Enterprise::VERSION_CENTER_SHOP_STORE, Enterprise::VERSION_DIRECT_CHAIN_STORE])) {
                $data['warehouse_type'] = EnterpriseConfig::MORE_WAREHOUSE;
            }
        }

        // 有门店不能为单仓
        if ($data['version'] != Enterprise::VERSION_CENTER_SHOP && $data['warehouse_type'] == EnterpriseConfig::SINGLE_WAREHOUSE) {
            error(201110);
        }

        // 企业唯一验证
        if (empty($data['corp_id'])) {
            error(201109);
        } else {
            if (Enterprise::query()->where('corp_id', '=', $data['corp_id'])->exists()) {
                notice(201007);
            }
        }

        if (Enterprise::query()->where('name', '=', $data['name'])->exists()) {
            notice(201006);
        }

        $enterpriseService         = new EnterpriseService();
        $enterpriseAdminService    = new EnterpriseAdminService();
        $enterpriseAdminNewService = new EnterpriseAdminNewService();

        if (empty($data['available_system'])) {
            $data['available_system'] = [1];
        }

        // 未传递服务商id则获取默认的服务商id
        if (!isset($data['service_merchant_id'])) {
            $defaultMerchant = (new Merchant())->defaultMerchant();
            if (empty($defaultMerchant)) {
                $merchantId = 6;// 默认写死
            } else {
                $merchantId = $defaultMerchant['id'];
            }
        } else {
            $merchantId = intval($data['service_merchant_id']);
        }
        // 企业信息添加
        try {
            DB::beginTransaction();

            $enterpriseData = [
                'hash'                => $enterpriseService->makeHash(),
                'name'                => $data['name'],
                'logo'                => $data['logo'] ?? '',
                'contact'             => $data['contact'] ?? '',
                'mobile'              => $data['mobile'] ?? '',
                'address'             => $data['address'] ?? '',
                'industry_id'         => $data['industry_id'] ?? 0,
                'separate_switch'     => 1,
                'lock_time'           => 0,
                'proportion'          => 0,
                'start_fee'           => 0,
                'crs_status'          => $data['crs_status'] ?? 1,
                'store_status'        => $data['store_status'] ?? 1,
                'contact_status'      => $data['contact_status'] ?? 2,
                'corp_id'             => trim($data['corp_id']) ?? '',
                'domain'              => '',
                'external_secret'     => '',
                'status'              => 1,
                'business_source'     => $data['business_source'] ?? 0,
                'version'             => $data['version'] ?? Enterprise::NORMAL_VERSION,
                'version_tag'         => $data['version_tag'] ?? '正式版',// 正式版｜试用版
                'service_merchant_id' => $merchantId,
                'account_expire_at'   => 7952313600, // 到期时间 2222-01-01 00:00:00
                'account_start_at'    => $data['account_start_at'] ?? 0,
                'account_limit'       => $data['account_limit'] ?? 10,
                'erp_status'          => 1,
                'crm_status'          => 1,
                'return_address'      => $data['return_address'] ?? '',
                'is_split_flow'       => 0,
                'available_system'    => Json::encode($data['available_system']),
                'edit_at'             => time(),
            ];

            if (!empty($enterpriseData['corp_id'])) {
                // 从uc查是否已经注册企业了 如果没注册则先创建
                try {
                    $res = (new \CrsSdk\Work\Enterprise())->detailByCorpId([
                        'corp_id' => $enterpriseData['corp_id']
                    ]);

                    Log::info('从UC查到的企业信息', $res);

                    $enterpriseData['domain'] = $res['epEnumber'];
                } catch (\Exception $exception) {
                    // uc创建企业
                    try {
                        $res = (new \CrsSdk\Work\Enterprise())->add($data);

                        Log::info('从UC创建完企业', $res);

                        $enterpriseData['domain'] = $res['epEnumber'];
                    } catch (\Exception $exception) {
                        Log::info('同步UC创建企业失败-' . $exception->getMessage(), $data);
                        // notice(710006);
                        throw new Error($exception->getCode(), $exception->getMessage());
                    }
                }
            }

            // 创建企业
            $enterprise    = Enterprise::query()->create($enterpriseData);
            $enterpriseArr = $enterprise->toArray();

            // 设置企业仓库类型
            EnterpriseConfigService::setWarehouseConfig($enterprise->id, $data['warehouse_type']);

            // 应用商店智慧导购
            EnterpriseWorkPlugin::query()->create([
                'enterprise_id' => $enterprise->id,
                'corp_id'       => $enterprise->corp_id,
                'identifier'    => EnterpriseWorkPlugin::CRS_APP_ASSISTANT
            ]);

            // 添加超管
            $adminData             = [
                'enterprise_id'          => $enterprise->id,
                'username'               => $data['email'],
                'email'                  => $data['email'],
                'account'                => $data['account'] ?? $data['email'],
                'password'               => '',
                'salt'                   => $enterpriseAdminService->makeSalt(),
                'global_roaming'         => '',
                'mobile'                 => $data['mobile'] ?? '',
                'role'                   => 1,
                'is_super'               => EnterpriseAdmin::IS_SUPER,
                'last_login_time'        => 0,
                'status'                 => 1,
                'member_id'              => $data['memUserId'] ?? '',
                'operator_id'            => 0,
                'admin_available_system' => Json::encode($data['available_system']),
                'range_type'             => 1,
            ];
            $password              = $originPwd = $data['password'];
            $adminData['password'] = md5(md5($password) . $adminData['salt']);

            Log::debug('添加超管数据', $adminData);
            try {
                $admin    = $enterpriseAdminNewService->addSuperAdmin($adminData);
                $adminArr = $admin->toArray();
            } catch (\Exception $e) {
                Log::error('添加超管异常', [
                    'message'   => $e->getMessage(),
                    'file'      => $e->getFile(),
                    'code'      => $e->getCode(),
                    'adminData' => $adminData
                ]);
                throw new Error($e->getCode(), $e->getMessage());
            }

            Log::debug('添加超管成功', $adminArr);

            // 同步门店创建企业为门店
            switchAppHost($enterprise->id);

            try {
                $paramsStore = [
                    'enterprise_id' => $enterprise->id,
                    'name'          => $data['name'],
                    'is_store'      => 2
                ];

                //(new Store())->add($paramsStore);

                /*
                // 添加默认仓库
                $res = (new Open())->addWarehouse([
                    'enterprise_id'      => $enterprise->id,
                    'warehouse_name'     => '默认仓库',
                    'name'               => '',
                    'province'           => '',
                    'city'               => '',
                    'area'               => '',
                    'address'            => '',
                    'longitude'          => '',
                    'latitude'           => '',
                    'return_name'        => '',
                    'return_postal_code' => '',
                    'return_province'    => '',
                    'return_city'        => '',
                    'return_area'        => '',
                    'return_address'     => '',
                    'warehouse_type'     => 2,
                    'source_id'          => '',
                    'status'             => 1,
                    'is_extract'         => 1,
                    'business_hours'     => [
                        'day'        => '1,2,3,4,5,6,7',
                        'start_hour' => '01',
                        'end_hour'   => '23',
                    ],
                ]);

                // 添加供货关系
                (new Open())->addSupplier([
                    'enterprise_id' => $enterprise->id,
                    'warehouse_id'  => $res['id'],
                    'source_id'     => $enterprise->id,
                ]);
                */
            } catch (\Exception $exception) {
                Log::info('同步企业到门店失败' . $exception->getMessage());
                notice(710005);
                //throw new Error($exception->getCode(), $exception->getMessage());
            }

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            Log::error('添加企业事务异常', [
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine(),
                'data'    => $data
            ]);
            throw new Error($e->getCode(), $e->getMessage());
        }

        // 企业注册后置数据延迟处理
        Log::debug('企业注册要入队');
        $registerJob = new EnterpriseRegisterJob($enterpriseArr, $adminArr, $originPwd, $data['permission_list'] ?? []);
        dispatch($registerJob);

        Log::debug('企业注册入队成功');

        $enterpriseArr['originPwd'] = $originPwd;
        $enterpriseArr['admin']     = $adminArr;

        return $enterpriseArr;
    }

    /**
     * 编辑企业
     *
     * @param $data
     * @return array|void
     * @throws \Exception
     */
    public function enterpriseEdit($data)
    {
        $result = Enterprise::query()->where('name', $data['name'])->where('id', '<>', $data['enterprise_id'])->first();

        // 企业名称已存在
        if ($result) {
            return notice(201006);
        }

        $enterprise = Enterprise::query()->where('id', '=', $data['enterprise_id'])->first();

        // uc更新企业信息
        $ucParams              = $data;
        $ucParams['qy_domain'] = $enterprise->domain;

        if (!empty($data['corp_id'])) {
            if (empty($ucParams['qy_domain'])) {
                // 从uc查是否已经注册企业了 如果没注册则先创建
                try {
                    $res = (new \CrsSdk\Work\Enterprise())->detailByCorpId([
                        'corp_id' => $data['corp_id']
                    ]);

                    Log::info('从UC查到的企业信息', $res);

                    $ucParams['qy_domain'] = $res['epEnumber'];
                } catch (\Exception $exception) {
                    // uc创建企业
                    try {
                        $res = (new \CrsSdk\Work\Enterprise())->add($data);

                        Log::info('从UC创建完企业', $res);

                        $ucParams['qy_domain'] = $res['epEnumber'];
                    } catch (\Exception $exception) {
                        Log::error('同步UC创建企业失败-' . $exception->getMessage(), $data);
                        // notice(710006);
                        throw new Error($exception->getCode(), $exception->getMessage());
                    }
                }
            }
            try {
                (new \CrsSdk\Work\Enterprise())->update($ucParams);
            } catch (\Exception $exception) {
                Log::info('更新UC企业信息失败--' . $exception->getMessage(), $data);
                throw new Error($exception->getCode(), $exception->getMessage());
            }
        }

        if (!isset($data['available_system'])) {
            $data['available_system'] = [];
        }

        // 编辑数据拼接
        $updateData = [
            'name'              => $data['name'] ?? '',
            'contact'           => $data['contact'] ?? '',
            'mobile'            => $data['mobile'] ?? '',
            'address'           => $data['address'] ?? '',
            'domain'            => $ucParams['qy_domain'],
            'corp_id'           => $data['corp_id'] ?? '',
            'business_source'   => $data['business_source'] ?? 0,
            'account_limit'     => $data['account_limit'] ?? 0,
            'account_start_at'  => (int)$data['account_start_at'],
            'account_expire_at' => (int)$data['account_expire_at'],
            'available_system'  => json_encode($data['available_system'], JSON_UNESCAPED_UNICODE),
            'edit_at'           => $data['edit_at'],
            'industry_id'       => $data['industry_id'],
            'version_tag'       => $data['version_tag'] ?? '正式版',// 正式版｜试用版
        ];

        Enterprise::query()->where('id', $data['enterprise_id'])->update($updateData);

        // 查询超管信息
        $enterpriseAdminInfo = EnterpriseAdmin::query()
            ->where('enterprise_id', $data['enterprise_id'])
            ->where('is_super', EnterpriseAdmin::IS_SUPER)
            ->first()
            ->toArray();
        // 需要更新的管理员信息
        $updateEnterpriseAdminData = [
            "admin_available_system"=> json_encode($data['available_system'], true)
        ];
        // 判断是否需要更新密码
        if (!empty($data['password'])) {
            $password = md5(md5($data['password']) . $enterpriseAdminInfo['salt']);
            $updateEnterpriseAdminData['password'] = $password;
        }
        EnterpriseAdmin::updateAvailableSystem($data['enterprise_id'], $updateEnterpriseAdminData);

        // 设置企业仓库类型
        EnterpriseConfigService::setWarehouseConfig($data['enterprise_id'], $data['warehouse_type']);

        $hashLockKey   = redisKey('ENTERPRISE_HASH_INFO', ['hash' => $enterprise->hash]);
        $domainLockKey = redisKey('ENTERPRISE_DOMAIN_INFO', ['domain' => $enterprise->domain]);
        $corpIdLockKey = redisKey('ENTERPRISE_CORP_ID_INFO', ['domain' => $enterprise->corp_id]);

        // 企业信息编辑 .删除缓存
        delRedis([
            RedisKey::REDIS_APP_NOVICE . $enterprise->id,
            RedisKey::REDIS_APP_GROUP . $enterprise->id,
            RedisKey::REDIS_APP_SECKILL . $enterprise->id,
            RedisKey::REDIS_ENTERPRISE_INFO . $enterprise->hash,
            $hashLockKey,
            $domainLockKey,
            $corpIdLockKey
        ]);

        return $this->getEnterpriseDetail($data['enterprise_id']);
    }

    /**
     * 获取企业详情
     *
     * @param $enterpriseId
     * @param bool $only true: 只获取企业基本信息 false：获取包含超管账号等在内的详细信息
     * @param array $field
     * @return array
     */
    public function getEnterpriseDetail($enterpriseId, $only = true, $field = ['*'])
    {
        if ($only) {
            $enterprise = Enterprise::query()->where('id', $enterpriseId)->first();
            return $enterprise ? $enterprise->toArray() : [];
        }

        $enterprise = Enterprise::query()
            ->select([
                'enterprise.*',
                'enterprise_admin.username',
                'enterprise_admin.account',
                'enterprise_admin.global_roaming',
                'enterprise_admin.email',
                'enterprise_admin.mobile as admin_mobile',
                'enterprise_admin.range_type',
                'enterprise_admin.admin_available_system',
            ])
            ->leftJoin('enterprise_admin', 'enterprise.id', '=', 'enterprise_id')
            ->where('enterprise.id', $enterpriseId)
            ->where('is_super', 1)
            ->first();
        return $enterprise ? $enterprise->toArray() : [];
    }

    public function syncApplication($enterprise, $enterpriseId, $appIds)
    {
        // 应用ID不为空
        if (!empty($appIds)) {
            $identity = [];
            foreach ($appIds as $appId) {
                // 查询应用ID是否存在
                $appData = EnterpriseApp::query()->where('id', $appId)->first(['id', 'identity']);
                // 构造应用唯一标识数组，同步信息使用
                array_push($identity, $appData->identity);

                if (empty($appData->id ?? '')) {
                    return notice(203003);
                }

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

                $_SERVER['enterprise_hash'] = $enterprise->is_split_flow == 1 ? $enterprise->hash : 'public';

                //                (new CrsSdk\Auth\Application())->syncApplication([
                //                    'enterprise_id' => $enterpriseId,
                //                    'applications'  => $identity
                //                ]);
            } catch (\Exception $e) {
                Log::info('同步企业应用失败：' . $e->getMessage());
            }
        }

        return true;
    }

    /**
     * 清除授权信息缓存
     *
     * @param int $enterpriseID
     * @return int
     */
    public function clearAuthInfoCache(int $enterpriseID)
    {
        $enterprise = Enterprise::query()->where('id', '=', $enterpriseID)->first();
        if ($enterprise !== null) {
            $redis   = redis();
            $idKey   = sprintf(RedisKey::AUTH_APP_INFO, $enterprise->getAttribute('id'));
            $hashKey = sprintf(RedisKey::AUTH_APP_INFO_HASH, $enterprise->getAttribute('hash'));

            $miniKey = RedisKey::MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID . $enterprise->getAttribute('id');
            return $redis->del([$idKey, $hashKey, $miniKey]);
        }

        return 0;
    }
}