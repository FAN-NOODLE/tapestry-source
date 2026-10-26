<?php
/*
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Services\OpenPlatform;

use App\Models\EnterpriseOpenPlatform;
use CrsSdk\Common\QCloudCos;
use CrsSdk\Common\Qiniu;
use CrsSdk\Shop\Version;
use App\Common\Constant;
use App\Exceptions\Error;
use App\Models\Enterprise;
use CrsSdk\Shop\AutoRelease;
use App\Jobs\SubmitMiniProgram;
use App\Services\BaseService;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use CrsSdk\Shop\TemplatePageVersion;
use App\Services\EnterpriseLogService;
use EasyWeChat\OpenPlatform\Application;
use App\Models\EnterpriseThirdPartTester;
use GuzzleHttp\Exception\GuzzleException;
use App\Models\EnterpriseBaseCodeLibrary;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseMiniProgramAuditLog;
use App\Models\EnterpriseMiniProgramAutoPack;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use EasyWeChat\Kernel\Exceptions\InvalidConfigException;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;

/**
 * 小程序相关的服务
 * Class MiniProgramService.
 */
class MiniProgramService extends BaseService
{
    /**
     * @var EnterpriseMiniProgramAuditLog
     */
    protected $miniProgramAuditLog;

    /**
     * @var EnterpriseMiniProgramCodeLibrary
     */
    protected $miniProgramCodeLibrary;

    /**
     * @var EnterpriseThirdPartAccount
     */
    protected $thirdPartAccount;

    /**
     * @var EnterpriseThirdPartTester
     */
    protected $thirdPartTester;

    /**
     * @var Application
     */
    public $openPlatform;

    /**
     * @var EnterpriseBaseCodeLibrary
     */
    protected $baseCodeLibrary;

    /**
     * @var Enterprise
     */
    protected $enterprise;

    /**
     * 状态 1：审核成功
     *      0：未发布  1：审核中  2：已发布  3：审核失败 7:审核成功
     */
    const AUDITSUCCESS = 2;


    /**
     * 状态 1：审核失败
     *      0：未发布  1：审核中  2：已发布  3：审核失败 7:审核成功
     */
    const AUDITERROR = 3;


    /**
     * 状态 1：审核失败
     *      0：未发布  1：审核中  2：已发布  3：审核失败 7:审核成功
     */
    const AUDIT_SUCCESS_NEW = 7;

    public function __construct(
        EnterpriseBaseCodeLibrary $baseCodeLibrary,
        EnterpriseMiniProgramAuditLog $miniProgramAuditLog,
        EnterpriseMiniProgramCodeLibrary $miniProgramCodeLibrary,
        EnterpriseThirdPartAccount $thirdPartAccount,
        EnterpriseThirdPartTester $thirdPartTester,
        Enterprise $enterprise
    )
    {
        $this->miniProgramAuditLog    = $miniProgramAuditLog;
        $this->miniProgramCodeLibrary = $miniProgramCodeLibrary;
        $this->baseCodeLibrary        = $baseCodeLibrary;
        $this->thirdPartAccount       = $thirdPartAccount;
        $this->thirdPartTester        = $thirdPartTester;
        $this->enterprise             = $enterprise;
        $this->openPlatform           = LaravelWechat::openPlatform('default');
    }

    /**
     * 提交小程序审核.
     * @param $appId
     * @param null $codeLibraryId
     * @param array|null $itemList
     * @return bool
     * [
     *                        {
     *                        "address":"index",
     *                        "tag":"学习 生活",
     *                        "first_class": "文娱",
     *                        "second_class": "资讯",
     *                        "first_id":1,
     *                        "second_id":2,
     *                        "title": "首页"
     *                        }
     *                        {
     *                        "address":"page/logs/logs",
     *                        "tag":"学习 工作",
     *                        "first_class": "教育",
     *                        "second_class": "学历教育",
     *                        "third_class": "高等",
     *                        "first_id":3,
     *                        "second_id":4,
     *                        "third_id":5,
     *                        "title": "日志"
     *                        }
     *                        ]
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function submitAudit($appId, $codeLibraryId = null, array $itemList = null)
    {
        Log::debug('提交小程序审核', [
            'appId'         => $appId,
            'itemList'      => $itemList,
            'codeLibraryId' => $codeLibraryId
        ]);

        // 获取小程序配置
        $config = $this->thirdPartAccount->where(['app_id' => $appId, 'status' => 1])->first();
        switchAppHost($config->enterprise_id);
        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        // 获取代码库
        if (!$codeLibraryId) {
            $codeLibrary = $this->miniProgramCodeLibrary
                ->where('account_id', $config->id)
                ->orderBy('id', 'DESC')
                ->first();
        } else {
            $codeLibrary = $this->miniProgramCodeLibrary->where('id', $codeLibraryId)->first();
        }
        if (!$codeLibrary) {
            throw new Error(100000, '未查到小程序代码版本');
        }

        if (!$itemList) {
            // 获取分类
            $categories = $this->getCategory($appId);

            if (empty($categories)) {
                throw new Error(100000, '未查到小程序分类');
            }
            $category = $categories[0]; // 获取第一个分类进行提交

            // 获取页面
            $pages = $this->getPage($appId);
            if (empty($categories)) {
                throw new Error(100000, '未查到小程序页面');
            }
            $page     = $pages[0]; // 获取首页
            $pageInfo = [
                'address' => $page,
                'title'   => '首页',
                'tag'     => '商城',
            ];

            // 组合成审核的列表页
            $itemList[] = array_merge($pageInfo, $category);
        }

        $feedbackInfo = null;
        if ($config->type == EnterpriseThirdPartAccount::TYPE_QY_WX) {
            $feedbackInfo = '仅在企业微信运行使用';
        }

        // 提交审核
        try {
            $res = $this->openPlatform
                ->miniProgram($appId, $config->refresh_token)
                ->code
                ->submitAudit($itemList, $feedbackInfo);

            if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
                EnterpriseLogService::error('提交小程序审核请求错误', [
                    'appId'         => $appId,
                    'refresh_token' => $config->refresh_token,
                    'itemList'      => $itemList,
                    'message'       => $res
                ]);

                throw new Error(100000, '提交微信小程序审核报错，错误码：' . json_encode($res));
//                if ($res['errcode'] != 85009 && $res['errcode'] != 85085) {
//                    EnterpriseLogService::error('提交小程序审核请求错误', [
//                        'appId'         => $appId,
//                        'refresh_token' => $config->refresh_token,
//                        'itemList'      => $itemList,
//                        'message'       => $res
//                    ]);
//
//                    throw new Error(100000, '提交微信小程序审核报错，错误码：' . json_encode($res));
//                } else {
//                    return false;
//                }
            }
        } catch (\Exception $e) {
            Log::error('提交小程序审核报错', [
                'appId'   => $appId,
                'message' => $e->getMessage(),
                'line'    => $e->getLine(),
                'file'    => $e->getFile()
            ]);

            throw new Error(100010, '提交微信小程序审核错误：' . $e->getMessage());
        }

        // 添加审核记录
        EnterpriseMiniProgramAuditLog::query()
            ->create([
                'code_library_id' => $codeLibrary->id,
                'audit_id'        => $res['auditid'],
                'commit'          => json_encode($itemList),
                'status'          => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['AUDITING'],
            ]);

        // 修改代码版本审核状态
        $codeLibrary->update([
            'audit_status' => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['AUDITING']
        ]);

        return true;
    }

    /**
     * 获取最近小程序的微信审核结果.
     * @param $appId
     * @return array
     * @throws Error
     * @throws InvalidConfigException
     */
    public function recentAudit($appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where(['app_id' => $appId, 'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE])
            ->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 查询小程序审核结果
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getLatestAuditStatus();

        Log::debug('获取最近小程序的微信审核结果：' . $appId, [
            'appId'  => $appId,
            'result' => $res
        ]);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            Log::error('获取最近小程序的微信审核结果错误', [
                'appId'   => $appId,
                'message' => $res
            ]);

            throw new Error(100000, '调用微信小程序审核请求错误，错误码：' . json_encode($res));
        }

        // 重新组织返回结果
        $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['PENDING'];
        switch ($res['status']) {
            case 0:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['AUDITED'];
                break;

            case 1:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['REFUSED'];
                break;

            case 2:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['AUDITING'];
                break;
        }

        return [
            'wx_audit_id' => $res['auditid'],
            'status'      => $auditStatus,
            'reason'      => $res['reason'] ?? '',
        ];
    }

    /**
     * 根据小程序审核id 获取微信小程序审核结果.
     * @param $appId
     * @param $auditId
     * @return array|void
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getAuditInfoByWXAuditId($appId, $auditId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount->where(['app_id' => $appId, 'status' => 1])->first();

        if (!$config) {
            return notice(300105);
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getAuditStatus($auditId);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('调用微信小程序审核请求错误', [
                'appId'   => $appId,
                'auditId' => $auditId,
                'message' => $res
            ]);

            throw new Error(100000, '调用微信小程序审核请求错误，错误码：' . json_encode($res));
        }

        // 重新组织返回结果
        $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['PENDING'];
        switch ($res['status']) {
            case 0:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['AUDITED'];
                $this->activeRecentAuditByWxAuditId($auditId);
                break;

            case 1:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['REFUSED'];
                $reason      = $res['reason'] ?? null;
                $this->negativeAuditByWxAuditId($auditId, $reason);
                break;

            case 2:
                $auditStatus = EnterpriseMiniProgramCodeLibrary::$AUDIT_STATUS['AUDITING'];
                break;
        }

        $result = [
            'audit_id' => $auditId,
            'status'   => $auditStatus,
            'reason'   => $res['reason'] ?? '',
        ];

        return $result;
    }

    /**
     * 根据微信审核ID, 添加成功的审核状态.
     * @param       $wxAuditId
     * @param array $log
     */
    public function activeRecentAuditByWxAuditId($wxAuditId, $log = [])
    {
        Log::debug('小程序审核成功回调', [
            'wxAuditId'           => $wxAuditId,
            'miniProgramAuditLog' => $log
        ]);

        // 获取审核日志
        if (!$log) {
            $log = $this->miniProgramAuditLog->where('audit_id', $wxAuditId)->first();
        }

        if (!$log) {
            throw new Error(100000, '该小程序审核日志不存在');
        }

        // 修改审核记录表
        $this->miniProgramAuditLog
            ->where('id', $log->id)
            ->update(['status' => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['AUDITED']]);

        // 修改代码版本库最近一次审核状态
        $this->miniProgramCodeLibrary
            ->where('id', $log->code_library_id)
            ->update(['audit_status' => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['AUDITED']]);

        //查询关联版本id修改shop版本状态失败
        $version = $this->miniProgramCodeLibrary::query()
            ->where('id', $log->code_library_id)
            ->first(["version_id", "account_id"]);
        // 审核成功通知
        if (!empty($version)) {
            Log::info("审核成功回调 ：开始通知shop 版本审核状态");
            try {
                $enterpriseInfo = $this->thirdPartAccount::query()
                    ->leftJoin('enterprise', 'enterprise_third_part_account.enterprise_id', '=', 'enterprise.id')
                    ->where('enterprise_third_part_account.id', $version->account_id)
                    ->where('enterprise_third_part_account.status', $this->thirdPartAccount::STATUS_ACTIVE)
                    ->where('enterprise_third_part_account.type', $this->thirdPartAccount::TYPE_SHOP)
                    ->first(['enterprise.hash']);
                Log::info('审核成功回调 ：企业信息 ' . json_encode($enterpriseInfo));
                if (!empty($enterpriseInfo)) {
                    switchAppHost($enterpriseInfo->hash);
                    $_SERVER['enterprise_hash'] = $enterpriseInfo->hash;
                }
                (new AutoRelease())->saveStatus([
                    'id'     => $version->version_id,
                    'status' => self::AUDIT_SUCCESS_NEW,
                ]);
            } catch (\Exception $e) {
                Log::info("审核成功回调 ：通知shop 版本审核状态 失败 reason:" . $e->getMessage());
                Log::error('审核成功回调，通知shop失败', [
                    'auditId' => $wxAuditId,
                    'message' => $e->getMessage(),
                    'file'    => $e->getFile(),
                    'line'    => $e->getLine()
                ]);
            }

        }
    }

    /**
     * 根据微信审核ID, 添加失败记录.
     * @param $wxAuditId
     * @param $reason
     * @throws Error
     */
    public function negativeAuditByWxAuditId($wxAuditId, $reason)
    {
        Log::info("审核失败回调 ：更新代码审核记录状态 对应代码库审核状态");
        EnterpriseLogService::error('小程序审核失败回调', [
            'wxAuditId' => $wxAuditId,
            'reason'    => $reason
        ]);

        // 获取审核日志
        $log = $this->miniProgramAuditLog->where('audit_id', $wxAuditId)->first();
        if (!$log) {
            throw new Error(100000, '该小程序审核日志不存在');
        }

        // 修改审核记录表
        $this->miniProgramAuditLog
            ->where('id', $log->id)
            ->update(['status' => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['AUDITED'], 'reason' => $reason]);

        // 修改代码版本库最近一次审核状态
        $this->miniProgramCodeLibrary
            ->where('id', $log->code_library_id)
            ->update(['audit_status' => EnterpriseMiniProgramAuditLog::$AUDIT_STATUS['REFUSED']]);

        //查询关联版本id修改shop版本状态失败
        $version = $this->miniProgramCodeLibrary::query()
            ->where('id', $log->code_library_id)
            ->first(["version_id", "account_id"]);
        // 审核失败通知
        if (!empty($version)) {
            Log::info("审核失败回调 ：开始通知shop 版本审核状态");
            try {
                $enterpriseInfo = $this->thirdPartAccount::query()
                    ->leftJoin('enterprise', 'enterprise_third_part_account.enterprise_id', '=', 'enterprise.id')
                    ->where('enterprise_third_part_account.id', $version->account_id)
                    ->where('enterprise_third_part_account.status', $this->thirdPartAccount::STATUS_ACTIVE)
                    ->where('enterprise_third_part_account.type', $this->thirdPartAccount::TYPE_SHOP)
                    ->first(['enterprise.hash']);
                Log::info('审核失败回调 ：企业信息 ' . json_encode($enterpriseInfo));
                if (!empty($enterpriseInfo)) {
                    $_SERVER['enterprise_hash'] = $enterpriseInfo->hash;
                    switchAppHost($enterpriseInfo->hash);
                }
                (new AutoRelease())->saveStatus([
                    'id'     => $version->version_id,
                    'status' => self::AUDITERROR,
                    'reason' => $reason
                ]);
            } catch (\Exception $e) {
                EnterpriseLogService::error('审核失败回调通知shop失败', [
                    'wxAuditId' => $wxAuditId,
                    'message'   => $e->getMessage(),
                    'line'      => $e->getLine(),
                    'file'      => $e->getFile()
                ]);
                Log::info("审核失败回调 ：通知shop 版本审核状态 失败 reason:" . $e->getMessage());
            }

        } else {
            throw new Error(100000, '关联shop版本id不存在' . $version->version_id);
        }
    }

    /**
     * 发布审核成功的最近的小程序.
     * @param $appId
     * @param null $wxAuditId
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function release($appId, $wxAuditId = null)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        switchAppHost($config->enterprise_id);

        Log::debug('发布审核成功的最近的小程序', [
            'appId'     => $appId,
            'wxAuditId' => $wxAuditId
        ]);

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 发布审核通过的小程序
        $res = $this->openPlatform
            ->miniProgram($config->app_id, $config->refresh_token)
            ->code
            ->release();

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            Log::error('微信错误：代码上线失败', [
                'appId'     => $appId,
                'wxAuditId' => $wxAuditId,
                'message'   => $res
            ]);

            throw new Error(100000, '代码上线失败，微信错误：' . json_encode($res));
        }

        // 发布代码，使其上线
        if ($wxAuditId) {
            try {
                DB::beginTransaction();
                // 获取审核日志
                $log = $this->miniProgramAuditLog->where('audit_id', $wxAuditId)->first();
                if (!$log) {
                    throw new Error(100000, '该小程序审核日志不存在');
                }

                // 下线其他版本
                $this->miniProgramCodeLibrary
                    ->where([
                        'account_id' => $config->id,
                        'using'      => true
                    ])->update(['using' => false]);

                // 待更新数据
                $updateData = [
                    'online_status' => EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS['ONLINE'],
                    'using'         => true
                ];

                // 查询分阶段发布数据
                $grayRelease = $this->miniProgramCodeLibrary->where('id', $log->code_library_id)->first();
                // 发布模式：分阶段发布
                if (EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY == $grayRelease->release_model) {
                    // 查询小程序分阶段发布详情
                    $grayReleasePlan = $this->getGrayReleasePlan($appId);
                    $grayReleasePlan = $grayReleasePlan['gray_release_plan'] ?? [];

                    // 全量发布
                    $updateData['release_model'] = EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_ALL;
                    if (!empty($grayReleasePlan['status'])) {
                        $updateData['grp_status'] = $grayReleasePlan['status'];
                    }
                }

                // 上线该版本
                $this->miniProgramCodeLibrary
                    ->where('id', $log->code_library_id)
                    ->update($updateData);

                //查询关联版本id修改shop版本状态失败
                $version = $this->miniProgramCodeLibrary
                    ->select("version_id")
                    ->where('id', $log->code_library_id)
                    ->first();
                // 上线通知
                if (!empty($version)) {
                    (new AutoRelease())->saveStatus([
                        'id'     => $version->version_id,
                        'status' => self::AUDITSUCCESS,
                    ]);
                } else {
                    Log::error('发布小程序未找到CodeLibrary的Version信息', ['audit_id' => $wxAuditId]);
                }

                DB::commit();
            } catch (\Exception $e) {
                DB::rollBack();
                Log::error('发布代码，修改DB-fail', [
                    'appId'     => $appId,
                    'wxAuditId' => $wxAuditId,
                    'message'   => $e->getMessage(),
                    'file'      => $e->getFile(),
                    'line'      => $e->getLine()
                ]);
                throw new Error(100000, $e->getMessage());
            }
            try {
                if (!empty($version)) {
                    Log::debug('发布小程序修改首页版本信息：' . $config->enterprise_id, [
                        'version_id' => $version->version_id
                    ]);
                    // 修改首页版本数据
                    (new TemplatePageVersion())->release([
                        'version_id' => $version->version_id
                    ]);
                }
            } catch (\Exception $e) {
                Log::error('发布小程序修改首页版本信息异常：' . $config->enterprise_id, [
                    'audit_id'   => $wxAuditId,
                    'version_id' => $version->version_id,
                    'message'    => $e->getMessage(),
                    'file'       => $e->getFile(),
                    'line'       => $e->getLine()
                ]);
            }
        }

        return $res;
    }

    /**
     * 修改线上小程序状态
     * @param $appId
     * @param $action
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function changeOnlineStatus($appId, $action)
    {
        EnterpriseLogService::debug('修改线上小程序状态', [
            'appId'  => $appId,
            'action' => $action
        ]);
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where(['app_id' => $appId, 'status' => 1])
            ->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        if (!key_exists($action, EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS) && $action !== 'PENDING') {
            throw new Error(100000, '小程序上线功能，暂未提供该访问状态');
        }

        $reflection = [
            EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS['OFFLINE'] => 'close',
            EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS['ONLINE']  => 'open',
        ];

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 修改上线小程序的访问状态
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->changeVisitStatus($reflection[$action]); // 管理线上小程序的状态，默认可访问

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：线上代码管理', [
                'appId'   => $appId,
                'message' => $res
            ]);

            throw new Error(100000, '线上代码管理失败，微信错误：' . json_encode($res));
        }

        // 调整线上版本小程序访问状态
        $info = $this->changeRecentReleaseCodeStatusByAccountId($config->id, $action);

        return $info;
    }

    /**
     * 更新小程序上线状态
     * @param $account_id
     * @param $status
     * @return mixed
     */
    protected function changeRecentReleaseCodeStatusByAccountId($account_id, $status)
    {
        $info = $this->miniProgramCodeLibrary
            ->where([
                'account_id'    => $account_id,
                'using'         => true,
                'online_status' => [
                    'online_status',
                    '<>',
                    EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS['PENDING']
                ]
            ])->latest()
            ->first();

        $info->online_status = EnterpriseMiniProgramCodeLibrary::$RELEASE_STATUS[$status];
        $info->update();

        return $info;
    }

    /**
     * 获取草稿箱列表.
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getDrafts()
    {
        // draft_list:
        // [
        //  {
        //      create_time: 1488965944,
        //      user_version: 'VVV',
        //      user_desc: 'AAS',
        //      draft_id: 0
        //   },
        //   {
        //      create_time: 1504790906,
        //      user_version: '11',
        //      user_desc: '111111',
        //      draft_id: 4
        //    }
        //  ]
        $res = $this->openPlatform->code_template->getDrafts();
        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：拉取代码草稿箱', [
                'message' => $res
            ]);
            throw new Error(100000, '拉取代码草稿箱失败，微信错误：' . json_encode($res));
        }

        return $res['draft_list'];
    }

    /**
     * 获取对应版本最新的草稿
     * @param $miniProgramTypeVersionPrefix
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getLatestDraft($miniProgramTypeVersionPrefix)
    {
        // {
        //    create_time: 1504790906,
        //    user_version: '11',
        //    user_desc: '111111',
        //    draft_id: 4
        // }
        $drafts = multi_collect($this->getDrafts());
        $draft  = $drafts->filter(function ($draft) use ($miniProgramTypeVersionPrefix) {
            return str_contains($draft['user_version'], $miniProgramTypeVersionPrefix);
        })->sortByDesc('create_time')->first();

        return $draft;
    }

    /**
     * 将指定草稿添加到模版.
     * @param $draftId
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function addDraftToTemplate($draftId)
    {
        $res = $this->openPlatform->code_template->createFromDraft($draftId);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：将指定草稿添加到模版', [
                'draftId' => $draftId,
                'message' => $res
            ]);

            throw new Error(100000, '添加到模板错误，微信错误：' . json_encode($res));
        }

        return true;
    }

    /**
     * 获取模版列表.
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getTemplates()
    {
        // template_list:
        // [
        //  {
        //      create_time: 1488965944,
        //      user_version: 'VVV',
        //      user_desc: 'AAS',
        //      template_id: 0
        //   },
        //   {
        //      create_time: 1504790906,
        //      user_version: '11',
        //      user_desc: '111111',
        //      template_id: 4
        //    }
        //  ]
        $res = $this->openPlatform->code_template->list();

        Log::info("小程序模板列表：" . json_encode($res));
        Log::info("开放平台配置：" . $this->openPlatform->config['app_id']);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：获取模板列表', $res);
            throw new Error(100000, '获取模板列表错误，微信错误：' . json_encode($res));
        }

        return $res['template_list'];
    }

    /**
     * 获取系统最新的模板
     * @param $miniProgramTypeVersionPrefix
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getLatestTemplate($miniProgramTypeVersionPrefix)
    {
        // {
        //    create_time: 1504790906,
        //    user_version: '11',
        //    user_desc: '111111',
        //    template_id: 4
        // }
        $templates = multi_collect($this->getTemplates());
        $template  = $templates->filter(function ($template) use ($miniProgramTypeVersionPrefix) {
            return str_contains($template['user_version'], $miniProgramTypeVersionPrefix);
        })->sortByDesc('create_time')->first();

        return $template;
    }

    /**
     * 将指定模版删除.
     * @param $templateId
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function deleteTemplate($templateId)
    {
        $res = $this->openPlatform->code_template->delete($templateId);
        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信报错：删除指定模板', [
                'templateId' => $templateId,
                'message'    => $res
            ]);
            throw new Error(100000, '删除模板错误，微信错误：' . json_encode($res));
        }

        return true;
    }

    /**
     * 设置小程序交互域名接口.
     * @param string $appId
     * @param string $action
     * @param array $requestDomain
     * @param array $wsRequestDomain
     * @param array $uploadDomain
     * @param array $downloadDomain
     * @param array $bizDomain
     * @param string $refreshToken
     * @return array
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function setDomain(
        string $appId,
        string $action,
        string $refreshToken = '',
        int $enterpriseId
    )
    {

        if (empty($enterpriseId)) {
            throw new Error(100000, '企业id不能为空');
        }

        // 获取当前企业的信息
        $enterprise = Enterprise::query()
            ->where('id', '=', $enterpriseId)
            ->first();
        if (empty($enterprise) || empty($enterprise->open_platform_id)) {
            EnterpriseLogService::error('企业或关联配置为空', [
                'enterpriseId' => $enterpriseId,
            ]);
            throw new Error(100000, '企业或关联配置为空');
        }
        // 获取配置信息
        $openPlatFormConfig = EnterpriseOpenPlatform::query()
            ->where('id', $enterprise->open_platform_id)
            ->first();

        $domainArr       = json_decode($openPlatFormConfig->domain_json, true);
        $requestDomain   = explode(',', $domainArr['request_domain']);
        $wsRequestDomain = explode(',', $domainArr['ws_request_domain']);
        $uploadDomain    = explode(',', $domainArr['upload_domain']);
        $downloadDomain  = explode(',', $domainArr['download_domain']);
        $bizDomain       = explode(',', $domainArr['biz_domain']);

        if (!$refreshToken) {
            // 获取小程序配置
            $config = $this->thirdPartAccount
                ->where('app_id', $appId)
                ->first();

            if (!$config) {
                throw new Error(100000, '未查到对应绑定的小程序');
            }

            $refreshToken = $config->refresh_token;
        }

        if (!in_array($action, ['add', 'delete', 'set', 'get'])) {
            throw new Error(100000, '小程序暂不提供该种设置域名方法');
        }

        $params = [
            'action'          => $action,
            'requestdomain'   => $requestDomain,
            'wsrequestdomain' => $wsRequestDomain,
            'uploaddomain'    => $uploadDomain,
            'downloaddomain'  => $downloadDomain,
        ];

        Log::info('修改小程序domain参数：' . var_export($params, true));

        // 切换企业开放平台配置
        changeOpenPlatformConfig($enterpriseId, $this->openPlatform);

        $res = $this->openPlatform->miniProgram($appId, $refreshToken)->domain->modify($params);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            Log::info('修改小程序domain参数：' . var_export($params, true));

            EnterpriseLogService::error('修改小程序服务器域名微信报错', [
                'appId'        => $appId,
                'refreshToken' => $refreshToken,
                'params'       => $params,
                'message'      => $res
            ]);
            throw new Error(100000, '修改小程序服务器域名错误，微信错误信息：' . json_encode($res));
        }
        $bizDomain = array_filter($bizDomain);
        Log::info('设置小程序webView域名:' . json_encode($bizDomain));
        if (!empty($bizDomain)) {

            $bizRes = $this->openPlatform->miniProgram($appId, $refreshToken)->domain->setWebviewDomain($bizDomain, $action);
            // 错误码 885001 处理，官方给出解释：系统失败，重试即可
            if ($bizRes['errcode'] == 885001) {
                $bizRes = $this->openPlatform->miniProgram($appId, $refreshToken)->domain->setWebviewDomain($bizDomain, $action);
            }
            if ($bizRes['errcode'] != 0 || $bizRes['errmsg'] != 'ok') {
                EnterpriseLogService::error('修改小程序webView域名微信错误', [
                    'message' => $bizRes,
                    'params'  => $bizDomain
                ]);
                throw new Error(100000, '修改小程序webView域名错误，微信错误信息：' . json_encode($bizRes) . '| 提交参数:' . json_encode($bizDomain));
            }
        }

        return [
            'request_domain'    => key_exists('requestdomain', $res) ? $res['requestdomain'] : '',
            'ws_request_domain' => key_exists('wsrequestdomain', $res) ? $res['wsrequestdomain'] : '',
            'upload_domain'     => key_exists('uploaddomain', $res) ? $res['uploaddomain'] : '',
            'download_domain'   => key_exists('downloaddomain', $res) ? $res['downloaddomain'] : '',
        ];
    }

    /**
     * 为小程序账号应用小程序模板
     * @param $templateId
     * @param $appId
     * @param string $extJson
     * @param string $version
     * @param string $description
     * @param int $versionId
     * @param int $onlineType
     * @param int $enterpriseId
     * @return mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function commit(
        $templateId,
        $appId,
        string $extJson,
        string $version,
        string $description,
        int $versionId,
        int $onlineType = EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_AUTOMATIC,
        $enterpriseId = 0
    )
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE
            ])->first();

        Log::info('config数据:' . var_export($config, true));

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 处理超时
        $requestNum = 1;
        do {
            try {
                // 应用小程序模板
                $res  = $this->openPlatform
                    ->miniProgram($appId, $config->refresh_token)
                    ->code
                    ->commit($templateId, $extJson, $version, $description);
                $flag = false;
            } catch (\Exception $e) {
                $requestNum++;
                $flag = true;
                if ($requestNum > 2) {
                    throw new Error(100000, '应用小程序模板错误，微信错误：微信响应超时超过3次');
                }
                sleep(10);
            }
        } while ($flag);


        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('应用小程序模板报错', [
                'appId'         => $appId,
                'templateId'    => $templateId,
                'refresh_token' => $config->refresh_token,
                'extJson'       => $extJson,
                'version'       => $version,
                'description'   => $description,
                'message'       => $res
            ]);
            throw new Error(100000, '应用小程序模板错误，微信错误：' . json_encode($res));
        }

        // tag和branch取提审时传给build的值
        $autoPackInfo = EnterpriseMiniProgramAutoPack::query()
            ->where('version_id', $versionId)
            ->when($enterpriseId, function ($query) use ($enterpriseId) {
                return $query->where('enterprise_id', $enterpriseId);
            })
            ->first(['version_id', 'task_config']);
        $tag          = $config->release_tag;
        $branch       = $config->release_branch;
        if (!empty($autoPackInfo)) {
            EnterpriseLogService::debug('存入小程序版本数据：auto_log', $autoPackInfo->toArray());
            $taskConfig = json_decode($autoPackInfo->task_config, true);
            if (!empty($taskConfig)) {
                $tag    = empty($taskConfig['config']['tag']) ? $tag : $taskConfig['config']['tag'];
                $branch = empty($taskConfig['branch']) ? $branch : $taskConfig['branch'];
            }
        }

        $data = [
            'account_id'     => $config->id,
            'version'        => $version,
            'version_id'     => $versionId,    //shop版本关联id
            'description'    => $description,
            'extra_info'     => $extJson,
            'template_id'    => $templateId,
            'online_type'    => $onlineType,
            'release_tag'    => $tag,
            'release_branch' => $branch,
        ];

        Log::info('存入小程序版本数据:' . var_export($data, true));
        EnterpriseLogService::debug('存入小程序版本数据', $data);

        // 保存小程序版本
        $info = $this->miniProgramCodeLibrary->create($data);

        return $info;
    }

    /**
     * 获取小程序体验二维码
     * @param string $appId
     * @return string
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getExperienceQrcode(string $appId)
    {
        // 获取小程序授权记录
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 已经存在测试二维码
        if ($config->tester_qrcode) {
            return $config->tester_qrcode;
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        // 不存在生成一个，并保存
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getQrCode();

        // 用于判断是否是成功返回的二维码图片
        $resType = $res->getHeader('Content-Type');
        if (in_array('image/jpeg', $resType)) {
            // 将生成的图片上传至腾讯云
            $imageUrl = (new QCloudCos())->uploadContent($res);
            if ($imageUrl) {
                $config->tester_qrcode = $imageUrl;
                $config->update();
            }
        } else {
            EnterpriseLogService::error('获取小程序体验二维码失败报错', [
                'appId'   => $appId,
                'message' => $res
            ]);
            throw new Error(100000, '获取小程序体验二维码失败，微信返回：' . json_encode($res));
        }

        return $config->tester_qrcode;
    }

    /**
     * 输出小程序体验二维码
     * @param string $appId
     * @throws InvalidConfigException
     * @throws GuzzleException
     */
    public function showExperienceQrcode(string $appId)
    {
        // 获取小程序授权记录
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 已经存在测试二维码
        if ($config->tester_qrcode) {
            header('Content-Type: image/jpeg');
            $data = file_get_contents($config->tester_qrcode);
            echo $data;
            exit;
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        // 不存在生成一个，并保存
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getQrCode();

        // 用于判断是否是成功返回的二维码图片
        $resType = $res->getHeader('Content-Type');
        if (in_array('image/jpeg', $resType)) {
            // 将生成的图片上传至七牛云
            $imageUrl = (new Qiniu())->object($res, 'jpeg');
            if ($imageUrl) {
                $config->tester_qrcode = $imageUrl;
                $config->update();
            }
        } else {
            EnterpriseLogService::error('获取小程序体验二维码失败', [
                'appId'   => $appId,
                'message' => $res
            ]);

            throw new Error(100000, '获取小程序体验二维码失败，微信返回：' . json_encode($res));
        }

        // 返回路径
        header('Content-Type: image/jpeg');
        $data = file_get_contents($config->tester_qrcode);
        echo $data;
        exit;
    }

    /**
     * 添加体验者.
     * @param $wechatId
     * @param $appId
     * @param array $remark
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function addTester($wechatId, $appId, $remark = [])
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => 1
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        // 添加体验者
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->tester
            ->bind($wechatId);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：添加体验者报错', [
                'wechatId' => $wechatId,
                'appId'    => $appId,
                'message'  => $res
            ]);
            throw new Error(100000, '添加体验者错误，微信错误：' . json_encode($res));
        }

        // 记录体验者
        $this->thirdPartTester->create([
            'third_part_account_id' => $config->id,
            'user_account_id'       => $wechatId,
            'name'                  => key_exists('name', $remark) ? $remark['name'] : str_random(5),
            'remark'                => key_exists('remark', $remark) ? $remark['remark'] : '',
        ]);

        return true;
    }

    /**
     * 解绑体验者.
     * @param $testerId
     * @param $appId
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function deleteTester($testerId, $appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => 1
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        $tester = $this->thirdPartTester
            ->where('id', $testerId)
            ->first();
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        // 解除体验者
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->tester
            ->unbind($tester->user_account_id);

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信错误：删除体检者报错', [
                'testerId' => $testerId,
                'appId'    => $appId,
                'message'  => $res
            ]);

            throw new Error(100000, '删除体验者错误，微信错误：' . json_encode($res));
        }

        // 删除体验者
        $this->thirdPartTester->delete($testerId);

        return true;
    }

    /**
     * 获取具体小程序的可选类目.
     * @param $appId
     * @return array
     *               [
     *               {
     *               "first_class":"工具",
     *               "second_class":"备忘录"，
     *               "first_id":1,
     *               "second_id":2,
     *               }
     *               {
     *               "first_class":"教育",
     *               "second_class":"学历教育",
     *               "third_class":"高等"
     *               "first_id":3,
     *               "second_id":4,
     *               "third_id":5,
     *               }
     *               ]
     * @throws Error
     * @throws InvalidConfigException
     */
    public function getCategory($appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => 1
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getCategory();

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信报错：获取小程序分类错误', [
                'appId'   => $appId,
                'message' => $res
            ]);

            throw new Error(100000, '获取小程序分类错误，微信错误：' . json_encode($res));
        }

        return $res['category_list'];
    }

    /**
     * 获取具体小程序应用代码的页面配置.
     * @param $appId
     * @return array
     *               [
     *               "index",
     *               "page\/list",
     *               "page\/detail"
     *               ]
     * @throws Error
     * @throws InvalidConfigException
     */
    public function getPage($appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => 1
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getPage();

        if ($res['errcode'] != 0 || $res['errmsg'] != 'ok') {
            EnterpriseLogService::error('微信报错：获取小程序页面配置错误', [
                'appId'   => $appId,
                'message' => $res
            ]);
            throw new Error(100000, '获取小程序页面配置错误，微信错误：' . json_encode($res));
        }

        return $res['page_list'];
    }

    /**
     * 提交微信审核全部流程.
     * @param $thirdPartAccount
     * @param bool $checked
     * @param int $onlineType
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function submitVersion(
        $thirdPartAccount,
        $checked = false,
        $onlineType = EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_AUTOMATIC
    )
    {
        Log::info('进入submitVersion');
        Log::info('thirdPartAccount:' . var_export($thirdPartAccount->toArray(), true));
        Log::info('checked.' . $checked);

        EnterpriseLogService::info('提交微信审核,进入submitVersion', [
            'thirdPartAccount' => $thirdPartAccount->toArray(),
            'checked'          => $checked
        ]);

        $enterprise = $this->enterprise
            ->where('id', $thirdPartAccount->enterprise_id)
            ->first();

        if (!$enterprise) {
            throw new Error(100000, '企业不存在！');
        }
        // 替换企业开放平台
        changeOpenPlatformConfig($thirdPartAccount->enterprise_id, $this->openPlatform);
        // 商城小程序
        if (EnterpriseThirdPartAccount::TYPE_SHOP === intval($thirdPartAccount->type)) {
            // 获取平台最新小程序模板
            $template = $this->getLatestTemplate(EnterpriseThirdPartAccount::TYPE_SHOP_SHORT);
        } else if (EnterpriseThirdPartAccount::TYPE_QY_WX == $thirdPartAccount->type) {
            $template = $this->getLatestTemplate(EnterpriseThirdPartAccount::TYPE_QY_WX_SHORT);
        } else {
            throw new  Error(100000, '小程序模板不存在！');
        }

        if (!$template) {
            throw new  Error(100000, '小程序模板不存在！');
        }

        Log::info('template:' . var_export($template, true));

        // 提交体验小程序
        if (config('app.env') == 'prod') {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host      = 'https://' . config('client.api_host.shop');
                $work_host = 'https://' . config('client.api_host.work');
            } else {
                $host = 'https://' . config('client.api_host.work');
            }
            $environment = 'release';
        } else if (config('app.env') == 'test') {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host      = config('client.api_url.shop');
                $work_host = config('client.api_url.work');
            } else {
                $host = config('client.api_url.work');
            }
            $environment = 'trial';
        } else {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host      = config('client.api_url.shop');
                $work_host = config('client.api_url.work');
            } else {
                $host = config('client.api_url.work');
            }

            $environment = 'trial';
        }

        //获取小程序console后台返回的版本信息
        $versionId = 0;
        if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
            switchAppHost($enterprise->id);
            $versionInfo = (new Version())->getVersionInfo(['enterprise_id' => $enterprise->id, 'version' => $template['user_version']]);
            if (!empty($versionInfo['data'])) {
                $versionId = $versionInfo['data']['id'];
            }
        }

        // 给微信的附加参数
        $extJson = [
            'extEnable' => true,
            'extAppid'  => $thirdPartAccount->app_id,
            'ext'       => [
                'host'            => $host,
                'work_host'       => isset($work_host) ? $work_host : '',
                'environment'     => $environment,
                'enterpriseHash'  => $enterprise->hash,
                'corpName'        => $enterprise->name,
                'miniProgramName' => $thirdPartAccount->name ?? '',
                'apiRoot'         => '/app/v1/',
            ],
        ];

        Log::info('附加参数:' . var_export($extJson, true));

        Log::info('开始设置可信域名:' . var_export($extJson, true));

        // 更新可信域名
        $this->setDomain($thirdPartAccount->app_id,
            'set',
            $thirdPartAccount->refresh_token,
            $thirdPartAccount->enterprise_id
        );

        $extJsonString = json_encode($extJson);
        Log::info('执行commit:');
        $this->commit(
            $template['template_id'],
            $thirdPartAccount->app_id,
            $extJsonString,
            $template['user_version'],
            $template['user_desc'],
            $versionId,
            $onlineType,
            $thirdPartAccount->enterprise_id
        );

        // 提交微信审核
        if ($checked) {
            Log::info('提交微信审核:' . $thirdPartAccount->app_id);
            $this->submitAudit($thirdPartAccount->app_id);
        }

        return true;
    }

    /**
     * 获取指定模板
     * @param $version
     * @return array|mixed
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function getTemplate($version)
    {
        $versionInfo = [];
        $allTemplate = $this->getTemplates();

        if (!$allTemplate) {
            throw new  Error(100000, '小程序模板不存在！');
        }

        foreach ($allTemplate as $item) {
            if ($item['user_version'] == $version) {
                $versionInfo = $item;
            }
        }

        return $versionInfo;
    }

    /**
     * 提审指定版本
     * @param string $appId
     * @param string $version
     * @param bool $checked
     * @param int $onlineType
     * @return bool
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function submitSpecifiedVersion(
        string $appId,
        string $version,
        bool $checked = false,
        int $onlineType = EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_AUTOMATIC
    )
    {
        // 小程序配置
        $thirdPartAccount = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$thirdPartAccount) {
            EnterpriseLogService::error('提审指定版本，未查到对应绑定的小程序', ['app_id' => $appId]);
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        changeOpenPlatformConfig($thirdPartAccount->enterprise_id, $this->openPlatform);
        Log::info('进入submitSpecifiedVersion');
        Log::info('version.' . $version);
        Log::info('thirdPartAccount:' . var_export($thirdPartAccount->toArray(), true));

        EnterpriseLogService::info('提审指定版本,进入submitSpecifiedVersion', [
            'version'          => $version,
            'thirdPartAccount' => $thirdPartAccount->toArray()
        ]);

        // 企业信息
        $enterprise = $this->enterprise->where('id', $thirdPartAccount->enterprise_id)->first();
        if (!$enterprise) {
            throw new Error(100000, '企业不存在！');
        }

        // 切换SHOP域名
        switchAppHost($enterprise->id);

        // 小程序版本对应的模板信息
        $template = $this->getTemplate($version);
        if (!$template) {
            throw new  Error(100000, '小程序指定模板' . $version . '不存在！');
        }

        Log::info('template:' . var_export($template, true));

        // 提交体验小程序
        if (config('app.env') == 'prod') {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host      = 'https://shop-api-crs.vchangyi.com';
                $work_host = 'https://work-api-crs.vchangyi.com';
            } else {
                $host = 'https://work-api-crs.vchangyi.com';
            }
            $environment = 'release';
        } else if (config('app.env') == 'test') {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host = 'https://crs-api-test.vchangyi.com';
                if (in_array($enterprise->hash, config('app.project_hash'))) {
                    $host = 'https://shop-api-test-crs.vchangyi.com';
                }
                $work_host = 'https://work-api-test-crs.vchangyi.com';
            } else {
                $host = 'https://work-api-test-crs.vchangyi.com';
            }
            $environment = 'trial';
        } else {
            if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
                $host      = 'https://shop-api-dev-crs.vchangyi.com';
                $work_host = 'https://work-api-dev-crs.vchangyi.com';
            } else {
                $host = 'https://work-api-dev-crs.vchangyi.com';
            }
            $environment = 'trial';
        }

        // 获取小程序console后台返回的版本信息
        $versionId = 0;
        if (EnterpriseThirdPartAccount::TYPE_SHOP == $thirdPartAccount->type) {
            $consoVersion = substr($version, 10);
            $versionInfo  = (new Version())->getVersionInfo(['enterprise_id' => $enterprise->id, 'version' => $consoVersion]);
            if (!empty($versionInfo['data'])) {
                $versionId = $versionInfo['data']['id'];
            }
        }

        // 给微信的附加参数
        $extJson = [
            'extEnable' => true,
            'extAppid'  => $thirdPartAccount->app_id,
            'ext'       => [
                'host'            => $host,
                'work_host'       => isset($work_host) ? $work_host : '',
                'environment'     => $environment,
                'enterpriseHash'  => $enterprise->hash,
                'corpName'        => $enterprise->name,
                'miniProgramName' => $thirdPartAccount->name ?? '',
                'apiRoot'         => '/app/v1/',
            ],
        ];

        Log::info('附加参数:' . var_export($extJson, true));

        Log::info('开始设置可信域名:' . var_export($extJson, true));

        // 更新可信域名
        $this->setDomain(
            $thirdPartAccount->app_id,
            'set',
            $thirdPartAccount->refresh_token,
            $thirdPartAccount->enterprise_id
        );

        $extJsonString = json_encode($extJson);
        Log::info('执行commit:');
        EnterpriseLogService::debug('执行commit', [
            $template['template_id'],
            $thirdPartAccount->app_id,
            $extJsonString,
            $template['user_version'],
            $template['user_desc'],
            $versionId,
            $onlineType
        ]);
        $this->commit(
            $template['template_id'],
            $thirdPartAccount->app_id,
            $extJsonString,
            $template['user_version'],
            $template['user_desc'],
            $versionId,
            $onlineType,
            $thirdPartAccount->enterprise_id
        );

        // 提交微信审核
        if ($checked) {
            Log::info('提交微信审核:' . $thirdPartAccount->app_id);
            $this->submitAudit($thirdPartAccount->app_id);
        }

        return true;
    }

    /**
     * 全网发布代码
     * @param $versionType
     * @param bool $checked
     * @param bool $is_query
     * @return array|string
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function fullProcedure($versionType, $checked = false, $is_query = false)
    {
        if (!in_array($versionType, EnterpriseThirdPartAccount::MINI_PROGRAM_VERSION)) {
            throw new Error(100000, '不存在该小程序');
        }

        $versionPrefix = $versionType;

        $type = EnterpriseThirdPartAccount::TYPE_QY_WX;
        if (EnterpriseThirdPartAccount::TYPE_SHOP_SHORT == $versionType) {
            $type = EnterpriseThirdPartAccount::TYPE_SHOP;
        }

        // 未到达走发布流程
        $thirdPartAccounts = $this->thirdPartAccount
            ->where([
                'platform' => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,
                'status'   => EnterpriseThirdPartAccount::STATUS_ACTIVE,
                'type'     => $type
            ])->get();

        $enterpriseIds = array_column($thirdPartAccounts->toArray(), 'enterprise_id');

        // 处理多开放平台
        $enterprise      = Enterprise::query()
            ->distinct()
            ->whereIn('id', $enterpriseIds)
            ->get(['open_platform_id'])
            ->groupBy('open_platform_id');
        $openPlatformIds = array_unique(array_keys($enterprise->toArray()));
        foreach ($openPlatformIds as $openPlatformId) {

            changeOpenPlatformConfig(0, $this->openPlatform, $openPlatformId);
            // 根据小程序代码类型，获取草稿箱最新的草稿
            $draft = $this->getLatestDraft($versionPrefix);

            // 将草稿箱的上传到模板箱
            $this->addDraftToTemplate($draft['draft_id']);

            // 根据小程序代码类型，获取模板箱列表
            $templates = $this->getTemplates();
            $templates = multi_collect($templates)->filter(function ($template) use ($versionPrefix) {
                return str_contains($template['user_version'], $versionPrefix);
            })->sortByDesc('create_time');

            // 添加版本到版本库
            if (count($templates)) {
                $template = $templates->first();

                $where = [
                    'template_id' => $template['template_id'],
                    'type'        => $versionType == 'q' ? 2 : 1
                ];

                $update = [
                    'version'     => $template['user_version'],
                    'description' => $template['user_desc']
                ];

                $query = $this->baseCodeLibrary->where($where)->first();
                if (!empty($query)) {
                    $query->update($update);
                } else {
                    $this->baseCodeLibrary->create($update);
                }
            }

            // 判断是否到达上限（3）
            if (count($templates) >= 3) {
                // 到达上线删除模板
                $lastTemplate = $templates->last();
                $this->deleteTemplate($lastTemplate['template_id']);
            }
        }


        // 使用队列
        if ($is_query) {
            dispatch(new SubmitMiniProgram($thirdPartAccounts, $checked));

            return 'success';
        } else {
            // 记录未发布的企业和错误
            $cannot = [];
            foreach ($thirdPartAccounts as $thirdPartAccount) {
                try {
                    $this->submitVersion($thirdPartAccount, $checked);
                } catch (\Exception $e) {
                    $cannot[] = [
                        'enterprise_id' => $thirdPartAccount->enterprise_id,
                        'account_id'    => $thirdPartAccount->id,
                        'error'         => [
                            'message' => $e->getMessage(),
                            'line'    => $e->getLine(),
                            'file'    => $e->getFile(),
                        ],
                    ];
                    EnterpriseLogService::error('全网发布代码失败', $cannot);
                    Log::info('错误数据:' . var_export($cannot, true));
                }
            }

            if (empty($cannot)) {
                return 'success';
            }

            return $cannot;
        }
    }

    /**
     * 撤回小程序审核
     * @param string $appId 小程序app_id
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws InvalidConfigException
     */
    public function cancelAudit(string $appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 小程序审核撤回
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->withdrawAudit();

        if (0 == $res['errcode'] || 'ok' == $res['errmsg']) {
            $codeLibrary = $this->miniProgramCodeLibrary
                ->where([
                    'account_id'   => $config->id,
                    'audit_status' => EnterpriseMiniProgramCodeLibrary::AUDITING
                ])->orderBy('id', 'DESC')
                ->first();

            if ($codeLibrary) {
                // 审核状态：撤审
                $codeLibrary->audit_status = EnterpriseMiniProgramCodeLibrary::WITHDRAW;
                $codeLibrary->update();

                $auditLog = $this->miniProgramAuditLog
                    ->where(['code_library_id' => $codeLibrary->id])
                    ->first();

                if ($auditLog) {
                    // 提交审核状态
                    $auditLog->status = EnterpriseMiniProgramCodeLibrary::WITHDRAW;
                    $auditLog->update();
                }
            }
        }

        return $res;
    }

    /**
     * 小程序列表数据
     * @param array $params 参数列表
     * @return mixed
     */
    public function getMiniPrograms(array $params = [])
    {
        // 关键字（小程序名称/企业名称/版本号）
        $keywords = isset($params['keywords']) ? $params['keywords'] : '';
        if ($keywords) {
            $keywords = strtr($keywords, [
                '%' => '\%',
                '_' => '\_'
            ]);
        }

        // 小程序类型
        $type = $params['type'] ?? 0;
        // 每页数据量
        $limit = $params['limit'] ?? Constant::DEFAULT_LIMIT;

        // 审核状态
        $auditStatus = $params['audit_status'] ?? '';
        if ($auditStatus) {
            $auditStatus = EnterpriseThirdPartAccount::AUDIT_STATUS[$auditStatus];
        }

        // 企业id集合
        $enterpriseIds = [];
        // 如果查询企业名称
        if ('' != $keywords) {
            // 企业id集合赋值
            $enterpriseIds = $this->enterprise
                ->where('name', 'LIKE', '%' . $keywords . '%')
                ->pluck('id')
                ->toArray();
        }

        // 查询条件
        $where = [
            'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 授权状态：有效
            'platform'      => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT, // 授权类型：小程序
            'verify_status' => EnterpriseThirdPartAccount::VERIFIED // 认证状态: 已认证
        ];

        // 小程序类型（1：商城，2：CRS）
        if ($type) {
            $where['type'] = $type;
        }

        // 查询关联的小程序提审记录最大id的sql语句
        $sql = "(select max(`id`) from enterprise_mini_program_code_library 
        where deleted_at=0 and account_id=enterprise_third_part_account.id)";

        return $this->thirdPartAccount
            ->select('id', 'app_id', 'enterprise_id', 'name', 'platform', 'type', 'status')
            ->where($where)
            ->where(function ($query) use ($keywords, $auditStatus, $enterpriseIds, $sql) {
                $query->when('' != $keywords, function ($query) use ($keywords, $sql) {
                    $query->where('name', 'LIKE', '%' . $keywords . '%')
                        ->orWhere(function ($query) use ($keywords) {
                            $query->whereHas('usedCodeLibrary', function ($query) use ($keywords) {
                                return $query->where('version', 'LIKE', '%' . $keywords . '%');
                            });
                        })->orWhere(function ($query) use ($keywords, $sql) {
                            $query->whereHas('latestCodeLibrary', function ($query) use ($keywords, $sql) {
                                return $query->where('version', 'LIKE', '%' . $keywords . '%')
                                    ->where('audit_status', EnterpriseMiniProgramCodeLibrary::AUDITING)
                                    ->where('id', DB::raw($sql));
                            });
                        });
                })->when($enterpriseIds, function ($query) use ($enterpriseIds) {
                    $query->orWhere(function ($query) use ($enterpriseIds) {
                        $query->whereIn('enterprise_id', $enterpriseIds);
                    });
                });
            })->when($auditStatus, function ($query) use ($auditStatus, $sql) {
                $query->whereHas('latestCodeLibrary', function ($query) use ($auditStatus, $sql) {
                    return $query->where([
                        'audit_status' => $auditStatus,
                        'id'           => DB::raw($sql)
                    ])->when($auditStatus == EnterpriseMiniProgramCodeLibrary::AUDITED, function ($query) use ($auditStatus) {
                        return $query->where([
                            'online_status' => EnterpriseMiniProgramCodeLibrary::PENDING
                        ]);
                    });
                });
            })->with(['usedCodeLibrary', 'latestCodeLibrary'])
            ->orderBy('id', 'DESC')
            ->paginate($limit)
            ->toArray();
    }


    /**
     * 小程序撤审次数（当天/当月）
     * @param array $accountIds 授权记录id集合
     * @return array
     */
    public function withdrawAuditCount(array $accountIds = [])
    {
        // 初始化待返数据
        $data = [];
        if (!$accountIds) {
            return $data;
        }

        // 当月开始时间
        $begin = strtotime(date('Y-m-1'));
        // 当前时间
        $end = time();
        // 当天开始时间
        $dayBegin = strtotime(date('Y-m-d', $end) . ' 0:0:0');

        // 当月开始时间-当前时间的撤审记录
        $codeLibraryData = $this->miniProgramCodeLibrary
            ->where('audit_status', EnterpriseMiniProgramCodeLibrary::WITHDRAW)
            ->whereIn('account_id', $accountIds)
            ->whereBetween('updated_at', [$begin, $end])
            ->get(['account_id', 'updated_at'])
            ->toArray();

        foreach ($codeLibraryData as $item) {
            // 授权记录id
            $accountId = $item['account_id'];
            // 撤审时间
            $updatedAt = $item['updated_at'];

            // 撤审时间大于等于当天开始时间 && 撤审时间小于等于当前时间
            if ($updatedAt >= $dayBegin && $updatedAt <= $end) {
                if (isset($data[$accountId]['day_count'])) {
                    $data[$accountId]['day_count']++;
                } else {
                    $data[$accountId]['day_count'] = 1;
                }
            }

            if (isset($data[$accountId]['month_count'])) {
                $data[$accountId]['month_count']++;
            } else {
                $data[$accountId]['month_count'] = 1;
            }
        }

        return $data;
    }

    /**
     * 小程序版本回退
     * @param string $appId 小程序app_id
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws InvalidConfigException
     */
    public function rollbackRelease(string $appId)
    {
        EnterpriseLogService::debug('小程序版本回退', [
            'app_id' => $appId
        ]);

        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        switchAppHost($config->enterprise_id);

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 小程序版本回退
        $res = $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->rollbackRelease();

        // 版本回退成功
        if (0 == $res['errcode'] || 'ok' == $res['errmsg']) {
            // 上一可回退版本
            $beforeCodeLibrary = $this->miniProgramCodeLibrary
                ->where([
                    'account_id'    => $config->id,
                    'audit_status'  => EnterpriseMiniProgramCodeLibrary::AUDITED, // 审核通过
                    'online_status' => EnterpriseMiniProgramCodeLibrary::ONLINE, // 上线
                    'using'         => EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE, // 是否在使用该版本：否
                    'release_model' => EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_ALL // 全量发布
                ])->orderByDesc('id')
                ->first();

            if ($beforeCodeLibrary) {
                // 当前在使用版本改为非当前使用版本
                $this->miniProgramCodeLibrary
                    ->where([
                        'account_id' => $config->id,
                        'using'      => EnterpriseMiniProgramCodeLibrary::USING_STATUS_TRUE
                    ])->update(['using' => EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE]);

                // 上一可回退版本改为当前使用版本
                $this->miniProgramCodeLibrary
                    ->where('id', $beforeCodeLibrary->id)
                    ->update(['using' => EnterpriseMiniProgramCodeLibrary::USING_STATUS_TRUE]);

                try {
                    EnterpriseLogService::debug('小程序版本回退修改首页版本信息：' . $config->enterprise_id, [
                        'version_id' => $beforeCodeLibrary->version_id
                    ]);
                    // 修改首页版本数据
                    (new TemplatePageVersion())->release(['version_id' => $beforeCodeLibrary->version_id]);
                } catch (\Exception $e) {
                    Log::error('小程序分阶段发布修改首页版本数据异常：' . $config->enterprise_id, [
                        'app_id'  => $appId,
                        'message' => $e->getMessage(),
                        'file'    => $e->getFile(),
                        'line'    => $e->getLine()
                    ]);
                    EnterpriseLogService::error('小程序分阶段发布修改首页版本数据异常：' . $config->enterprise_id, [
                        'app_id'  => $appId,
                        'message' => $e->getMessage(),
                        'file'    => $e->getFile(),
                        'line'    => $e->getLine()
                    ]);
                }
            }
        }

        return $res;
    }

    /**
     * 获取小程序分阶段发布数据
     * @param array $accountIds
     * @return array
     */
    public function getGrayReleaseData(array $accountIds = [])
    {
        // 初始化待返数据
        $data = [];
        if (!$accountIds) {
            return $data;
        }

        // 分阶段发布状态：执行完成、被删除
        $releaseStatus = [
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_COMPLETED,
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_DELETED,
        ];

        // 查询小程序分阶段发布数据
        return $this->miniProgramCodeLibrary
            ->where([
                'release_model' => EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY,
                'using'         => EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE
            ])->whereIn('account_id', $accountIds)
            ->whereNotIn('grp_status', $releaseStatus)
            ->get()
            ->toArray();
    }

    /**
     * 小程序分阶段发布详情
     * @param string $appId
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws InvalidConfigException
     */
    public function getGrayReleasePlan(string $appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 小程序分阶段发布详情
        return $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->getGrayRelease();
    }

    /**
     * 更新小程序分阶段发布数据
     * @param int $id 主键id
     * @param array $data 待更新数据
     * @return mixed
     */
    public function updateGrayReleaseStatus(int $id, array $data)
    {
        return $this->miniProgramCodeLibrary->where('id', $id)->update($data);
    }

    /**
     * 小程序分阶段发布
     * @param string $appId
     * @param int $grayPercentage
     * @param null $wxAuditId
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function grayRelease(string $appId, int $grayPercentage, $wxAuditId = null)
    {
        EnterpriseLogService::debug('小程序分阶段发布', [
            'app_id'              => $appId,
            'grp_gray_percentage' => $grayPercentage,
            'audit_id'            => $wxAuditId
        ]);

        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        switchAppHost($config->enterprise_id);

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 发布审核通过的小程序
        $res = $this->openPlatform
            ->miniProgram($config->app_id, $config->refresh_token)
            ->code
            ->grayRelease($grayPercentage);

        if (0 != $res['errcode'] || 'ok' != $res['errmsg']) {
            EnterpriseLogService::error('微信报错：小程序分阶段发布失败', [
                'app_id'         => $appId,
                'grayPercentage' => $grayPercentage,
                'message'        => $res
            ]);

            throw new Error(100000, '小程序分阶段发布失败，微信错误：' . json_encode($res));
        }

        // 发布代码
        if ($wxAuditId) {
            // 获取审核日志
            $log = $this->miniProgramAuditLog->where('audit_id', $wxAuditId)->first();
            if (!$log) {
                throw new Error(100000, '该小程序审核日志不存在');
            }

            // 查询小程序分阶段发布详情
            $grayReleasePlan = $this->getGrayReleasePlan($appId);
            $grayReleasePlan = $grayReleasePlan['gray_release_plan'] ?? [];

            // 更新小程序版本数据
            $this->miniProgramCodeLibrary
                ->where('id', $log->code_library_id)
                ->update([
                    'release_model'       => EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY, // 分阶段发布
                    'grp_gray_percentage' => $grayPercentage,
                    'grp_status'          => !empty($grayReleasePlan['status']) ? $grayReleasePlan['status'] : 0
                ]);

            try {
                $version = $this->miniProgramCodeLibrary
                    ->select('version_id')
                    ->where('id', $log->code_library_id)
                    ->first();
                EnterpriseLogService::error('小程序分阶段发布修改首页版本数据：' . $config->enterprise_id, [
                    'app_id'              => $appId,
                    'grp_gray_percentage' => $grayPercentage,
                    'audit_id'            => $wxAuditId,
                    'version_id'          => $version->version_id
                ]);
                // 修改首页版本数据
                (new TemplatePageVersion())->release([
                    'version_id' => $version->version_id
                ]);
            } catch (\Exception $e) {
                Log::error('小程序分阶段发布修改首页版本数据异常：' . $config->enterprise_id, [
                    'app_id'              => $appId,
                    'grp_gray_percentage' => $grayPercentage,
                    'audit_id'            => $wxAuditId,
                    'message'             => $e->getMessage(),
                    'file'                => $e->getFile(),
                    'line'                => $e->getLine()
                ]);
                EnterpriseLogService::error('小程序分阶段发布修改首页版本数据异常：' . $config->enterprise_id, [
                    'app_id'              => $appId,
                    'grp_gray_percentage' => $grayPercentage,
                    'audit_id'            => $wxAuditId,
                    'message'             => $e->getMessage(),
                    'file'                => $e->getFile(),
                    'line'                => $e->getLine()
                ]);
            }

        }

        return $res;
    }

    /**
     * 小程序取消分阶段发布
     * @param string $appId
     * @return array|\EasyWeChat\Kernel\Support\Collection|object|\Psr\Http\Message\ResponseInterface|string
     * @throws InvalidConfigException
     */
    public function revertGrayRelease(string $appId)
    {
        // 获取小程序配置
        $config = $this->thirdPartAccount
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }

        // 切换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        // 发布审核通过的小程序
        $res = $this->openPlatform
            ->miniProgram($config->app_id, $config->refresh_token)
            ->code
            ->revertGrayRelease();

        if (0 != $res['errcode'] || 'ok' != $res['errmsg']) {
            EnterpriseLogService::error('微信报错：小程序取消分阶段发布失败', [
                'appId'   => $appId,
                'message' => $res
            ]);
            throw new Error(100000, '小程序取消分阶段发布失败，微信错误：' . json_encode($res));
        }

        // 分阶段发布记录
        $codeLibrary = $this->miniProgramCodeLibrary
            ->where([
                'account_id'    => $config->id,
                'release_model' => EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY,
                ['grp_status', '<>', EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_DELETED]
            ])->orderBy('id', 'DESC')
            ->first();

        if ($codeLibrary) {
            // 状态：被删除
            $codeLibrary->grp_status = EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_DELETED;
            $codeLibrary->update();
        }

        return $res;
    }

    /**
     * 提交指定版本号草稿 至 模版库
     * @param $version
     * @throws GuzzleException
     * @throws InvalidConfigException
     */
    public function addOneDraftToTemplate($version)
    {
        // 获取草稿箱列表
        $draftList = $this->getDrafts();
        $draft     = [];
        // 提交指定版本草稿至模版库
        foreach ($draftList as $item) {
            if ($item['user_version'] == $version) {
                $draft = $item;
            }
        }

        if (empty($draft)) {
            throw  new  Error(10000, '草稿不存在');
        }
        $this->addDraftToTemplate($draft['draft_id']);
    }

    /**
     * 根据商户版本id获取微信审核id
     * @param $versionId
     *
     * @return array
     */
    public function getAuditIdByVersionId($versionId, $enterpriseId = 0)
    {
        $accountId = 0;
        if ($enterpriseId) {
            $thirdPartAccount = $this->getThirdPartAccount($enterpriseId);
            if ($thirdPartAccount) {
                $accountId = $thirdPartAccount->id;
            }
        }
        $codeLibrary = EnterpriseMiniProgramCodeLibrary::query()
            ->where('version_id', $versionId)
            ->when($accountId, function ($query) use ($accountId) {
                return $query->where('account_id', $accountId);
            })
            ->where('audit_status', EnterpriseMiniProgramCodeLibrary::AUDITED)
            ->first(['id']);
        if (empty($codeLibrary)) {
            throw new Error(100000, '该小程序提审记录不存在');
        }
        $codeLibrary = $codeLibrary->toArray();
        $wxAuditInfo = EnterpriseMiniProgramAuditLog::query()
            ->where('code_library_id', $codeLibrary['id'])
            ->first(['audit_id']);
        if (empty($wxAuditInfo)) {
            throw new Error(100000, '该小程序审核日志不存在');
        }
        return $wxAuditInfo->toArray();
    }

    /**
     * 根据小程序app_id和版本获取审核id
     * @param $version
     * @param $appId
     *
     * @return array
     */
    public function getAuditIdByVersionAppId($version, $appId)
    {
        $config = EnterpriseThirdPartAccount::query()
            ->where(['app_id' => $appId, 'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE])
            ->first();
        if (!$config) {
            throw new Error(100000, '未查到对应绑定的小程序');
        }
        $codeLibrary = EnterpriseMiniProgramCodeLibrary::query()
            ->where('version', $version)
            ->where('account_id', $config->id)
            ->where('audit_status', EnterpriseMiniProgramCodeLibrary::AUDITED)
            ->first(['id']);
        if (empty($codeLibrary)) {
            throw new Error(100000, '该小程序提审记录不存在');
        }
        $wxAuditInfo = EnterpriseMiniProgramAuditLog::query()
            ->where('code_library_id', $codeLibrary->id)
            ->first(['audit_id']);
        if (empty($wxAuditInfo)) {
            throw new Error(100000, '该小程序审核日志不存在');
        }
        return ['wx_audit_id' => $wxAuditInfo->audit_id];
    }

    public function queryQuota($appId)
    {
        $config = $this->thirdPartAccount->where(['app_id' => $appId, 'status' => 1])->first();
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        return $this->openPlatform
            ->miniProgram($appId, $config->refresh_token)
            ->code
            ->queryQuota();
    }

    /**
     * @param $enterpriseId
     * @return mixed
     */
    public function getThirdPartAccount($enterpriseId)
    {
        return EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('verify_status', EnterpriseThirdPartAccount::VERIFIED)
            ->first(['id']);
    }
}
