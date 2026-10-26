<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:20
 */

namespace App\Services;

use App\Common\Constant;
use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Models\Enterprise;
use App\Cache\EnterpriseCache;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseApp;
use App\Models\EnterpriseAppInstall;
use CrsSdk\Console\Application;
use CrsSdk\Shop\PagePostage;
use CrsSdk\Shop\UserLevel;
use CrsSdk\Work\CommissionSetting;
use Exception;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use App\Http\Resources\EnterpriseResource;

class EnterpriseService extends BaseService
{
    /**
     * 已过期
     */
    const EXPIRED = 1;
    /**
     * 效期中
     */
    const NO_EXPIRE = 2;
    /**
     * 即将过期
     */
    const WILL_EXPIRED = 3;

    /**
     * 企业模型
     *
     * @var Enterprise
     */
    protected $enterprise;

    /**
     * 企业列表分页查询
     *
     * @param int $limit 每页显示条数
     * @param array $where 查询条件
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function paginate($limit, $where)
    {
        $query = Enterprise::query();

        // 条件组合
        if (!empty($where)) {
            foreach ($where as $key => $value) {
                if (!empty($value) || $value == 0) {
                    switch ($key) {
                        case 'name':
                            $query = $query->where($key, 'like', '%' . $value . '%');
                            break;
                        default:
                            $query = $query->where($key, $value);
                    }
                }
            }
        }

        // 排序
        $query = $query->orderBy('id', 'desc');

        try {
            return EnterpriseResource::collection($query->paginate((int)$limit));
        } catch (Exception $e) {
            Log::error('获取企业列表失败：' . $e->getMessage());
            return error(110000);
        }
    }

    /**
     * 根据企业ID获取企业详情
     *
     * @param int $enterpriseId 企业ID
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function enterprise($enterpriseId)
    {
        try {
            return Enterprise::query()->where('id', $enterpriseId)->first();
        } catch (Exception $e) {
            return error(110000);
        }
    }

    /**
     * 根据企业HASH获取企业详情
     *
     * @param string $hash 企业HASH
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function enterpriseHash($hash)
    {
        try {
            return Enterprise::query()->where('hash', $hash)->first();
        } catch (Exception $e) {
            return error(110000);
        }
    }

    /**
     * 添加企业信息
     *
     * @param array $data 添加的数据
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function enterpriseAdd($data)
    {
        $result = Enterprise::query()->where('name', $data['name'])->first();

        // 企业名称已存在
        if ($result) {
            return notice(201006);
        }

        try {
            if ($data['crs_status'] == Constant::ENTERPRISE_INSTALL_STATUS && !empty($data['corp_id'])) {
                $data['corp_id'] = trim($data['corp_id']);
                // uc创建企业
                $res            = (new \CrsSdk\Work\Enterprise())->add($data);
                $data['domain'] = $res['epEnumber'];
            }
            Log::info('创建企业数据', $data);
            return Enterprise::query()->create($data);
        } catch (Exception $e) {
            return error(110000);
        }
    }

    /**
     * 编辑企业信息
     *
     * @param int $enterpriseId 修改的企业ID
     * @param array $data 修改的信息
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function enterpriseSave($enterpriseId, $data)
    {
        // 排除重名
        $res = Enterprise::query()->where('name', $data['name'])->where('id', '<>', $enterpriseId)->first();

        if ($res) {
            return error(201006);
        }

        try {
            if ($data['crs_status'] == Constant::ENTERPRISE_INSTALL_STATUS && !empty($data['corp_id'])) {
                $data['corp_id'] = trim($data['corp_id']);
                $enterprise      = Enterprise::query()->where('id', $enterpriseId)->first();
                if (empty($enterprise->corp_id)) {
                    // uc创建企业
                    $ucRes          = (new \CrsSdk\Work\Enterprise())->add($data);
                    $data['domain'] = $ucRes['epEnumber'];
                } else {
                    // uc更新企业信息
                    $ucParams              = $data;
                    $ucParams['qy_domain'] = $enterprise->domain;
                    $ucRes                 = (new \CrsSdk\Work\Enterprise())->update($ucParams);
                    $data['domain']        = $ucRes['epEnumber'];
                }
            }

            $result = Enterprise::query()->where('id', $enterpriseId)->update($data);

            if ($result) {
                $enterprise      = $this->enterprise($enterpriseId);
                $enterpriseCache = new EnterpriseCache();
                $enterpriseCache->update()->detail('id', (string)$enterprise->id);
                $enterpriseCache->update()->detail('hash', $enterprise->hash);

                try {
                    // 更新企业名字给open
                    (new \CrsSdk\Open\Enterprise())->update([
                        'name'            => $enterprise->name,
                        'enterprise_hash' => $enterprise->hash
                    ]);
                } catch (Exception $e) {
                    Log::critical('更新企业信息给OPEN', [
                        'line'    => $e->getLine(),
                        'code'    => $e->getCode(),
                        'file'    => $e->getFile(),
                        'message' => $e->getMessage(),
                    ]);
                }
            }

            return $result;
        } catch (Exception $e) {
            return error(110000);
        }
    }

    /**
     * 创建不重复的HASH
     *
     * @return bool|string
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/28 11:16
     */
    public function makeHash()
    {
        $hash = md5(uniqid());
        if (Enterprise::query()->where('hash', $hash)->first()) {
            return $this->makeHash();
        }

        return $hash;
    }

    /**
     * 获取企业列表（新）
     *
     * @param array $params
     * @param int $page
     * @param int $perPage
     * @return LengthAwarePaginator
     */
    public function enterpriseList(array $params, int $page, int $perPage)
    {
        $enterpriseQuery = Enterprise::query()
            ->where('version', '>', Enterprise::NORMAL_VERSION)
            ->when(!empty($params['name']), function ($query) use ($params) {
                return $query->where('name', 'like', '%' . $params['name'] . '%');
            })->when(!empty($params['version']), function ($query) use ($params) {
                return $query->where('version', '=', $params['version']);
            })->when(!empty($params['version_tag']), function ($query) use ($params) {
                return $query->where('version_tag', '=', $params['version_tag']);
            })->when(isset($params['business_source']) && $params['business_source'] !== '', function ($query) use ($params) {
                return $query->where('business_source', '=', $params['business_source']);
            })->when(!empty($params['industry_id']), function ($query) use ($params) {
                return $query->where('industry_id', '=', $params['industry_id']);
            });
        // 已过期
        if (isset($params['expire_type']) && $params['expire_type'] == self::EXPIRED) {
            $enterpriseQuery = $enterpriseQuery->where('account_expire_at', '<', time());
        }
        // 未过期
        if (isset($params['expire_type']) && $params['expire_type'] == self::NO_EXPIRE) {
            $enterpriseQuery = $enterpriseQuery->where('account_expire_at', '>', time() + 2592000);
        }
        // 即将过期
        if (isset($params['expire_type']) && $params['expire_type'] == self::WILL_EXPIRED) {
            $enterpriseQuery = $enterpriseQuery->where('account_expire_at', '>', time())
                ->where('account_expire_at', '<', time() + 2592000);
        }

        return $enterpriseQuery->orderBy('id', 'DESC')
            ->paginate($perPage, ['*'], 'page', $page)
            ->toArray();
    }

    /**
     * 新增企业
     *
     * @param array $data
     * @param string $account 账号
     * @return mixed
     * @throws Exception
     */
    public function add($data, $account)
    {
        try {
            DB::beginTransaction();

            $enterprise = $this->enterpriseAdd([
                'name'           => $data['name'],
                'contact'        => $data['contact'],
                'mobile'         => $data['mobile'],
                'address'        => $data['address'],
                'crs_status'     => $data['crs_status'],
                'store_status'   => $data['store_status'],
                'corp_id'        => $data['corp_id'],
                'hash'           => $this->makeHash(),
                'version'        => $data['version'] ?? Enterprise::NORMAL_VERSION,
                'version_tag'    => $data['version_tag'] ?? '正式版',// 正式版｜试用版
                'edit_at'        => time(),
                'contact_status' => 2,// 通讯录默认禁用状态
            ]);
            if ($enterprise !== null) {

                // 切换SHOP域名
                switchAppHost($enterprise->id);

                // 同步企业名字和hash给Open配置表
                try {
                    (new \CrsSdk\Open\Enterprise())->add([
                        'name'            => $enterprise->name,
                        'enterprise_hash' => $enterprise->hash
                    ]);
                } catch (Exception $e) {
                    Log::error('创建企业:同步企业信息:open接口异常', [
                        'line'    => $e->getLine(),
                        'code'    => $e->getCode(),
                        'file'    => $e->getFile(),
                        'message' => $e->getMessage(),
                    ]);
                }

                // 生成默认的邮费配置
                try {
                    (new PagePostage())->save([
                        'enterprise_id' => $enterprise->id,
                        'on_off'        => 1,   // 1：开，2：关
                        'rule'          => 9900
                    ]);
                } catch (Exception $e) {
                    EnterpriseLogService::error('创建企业:生成默认邮费配置:shop接口异常', [
                        'line'    => $e->getLine(),
                        'code'    => $e->getCode(),
                        'file'    => $e->getFile(),
                        'message' => $e->getMessage(),
                    ]);
                }

                // 企业会员等级初始化
                try {
                    (new UserLevel())->initLevel(['enterprise_id' => $enterprise->id]);
                } catch (Exception $e) {
                    Log::info('企业添加会员等级初始化失败：' . $e->getMessage());
                }

                // 企业导购佣金设置初始化
                try {
                    (new CommissionSetting())->initCommission(['enterprise_id' => $enterprise->id]);
                } catch (Exception $e) {
                    Log::info('企业导购佣金设置初始化失败：' . $e->getMessage());
                }
            }
            $adminService = new EnterpriseAdminService();

            // 初始化超管
            $adminData             = [
                'enterprise_id' => $enterprise->id,
                'account'       => $account,
                'salt'          => $adminService->makeSalt(),
            ];
            $adminData['password'] = md5(md5(EnterpriseAdminService::PASSWORD) . $adminData['salt']);
            $adminService->adminAdd($adminData);

            // 安装应用
            $_SERVER['enterprise_hash'] = $enterprise->is_split_flow == 1 ? $enterprise->hash : 'public';
            $this->initApp($enterprise->id, $data['app_ids'] ?? []);

            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();

            throw new Exception($exception->getMessage());
        }

        return $enterprise;
    }

    /**
     * 初始化应用
     *
     * @param int $enterpriseId
     * @param array $appIds
     * @return bool|void
     * @throws Exception
     */
    private function initApp($enterpriseId, $appIds)
    {
        if (empty($appIds)) {
            return true;
        }

        // 首先初始化数据
        EnterpriseAppInstall::query()->where('enterprise_id', '=', $enterpriseId)->delete();

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
            } catch (Exception $e) {
                return notice(203004);
            }
        }

        // 开始同步企业应用到 console
        try {
            // 老版
            (new Application())->syncApplication([
                'enterprise_id' => $enterpriseId,
                'applications'  => $identity
            ]);
            // 新版（gw-console）
            (new Application())->gwSyncApplication([
                'enterprise_id' => $enterpriseId,
                'applications'  => $identity
            ]);
        } catch (Exception $e) {
            Log::info('同步企业应用失败：' . $e->getMessage());
        }
    }
}
