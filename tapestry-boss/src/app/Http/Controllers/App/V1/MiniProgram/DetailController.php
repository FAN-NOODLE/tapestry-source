<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-05-07
 * Time: 11:54
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterpriseMiniProgramAuditLog;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 小程序详情
     *
     * @param Request $request
     *
     * @return array|void
     */
    public function index(Request $request, MiniProgramService $miniProgramService)
    {
        // 小程序app_id
        $appId = $request->post('app_id', '');
        if (!$appId) {
            return error(300100);
        }

        // 小程序配置
        $config = EnterpriseThirdPartAccount::query()
            ->select(['id', 'app_id', 'name', 'type', 'enterprise_id'])
            ->where([
                'app_id' => $appId,
                'status' => EnterpriseThirdPartAccount::STATUS_ACTIVE // 授权状态：有效
            ])->first();

        if (!$config) {
            return error(300105);
        }

        // 替换企业开放平台信息
        changeOpenPlatformConfig($config->enterprise_id, $miniProgramService->openPlatform);

        // 获取小程序撤审次数（当天/当月）
        $withdrawAuditCount = $miniProgramService->withdrawAuditCount([$config->id]);
        // 初始化待返数据
        $data = $this->format($config, $withdrawAuditCount);

        return success(['data' => $data]);
    }

    /**
     * 格式化数据
     *
     * @param $config
     * @param $withdrawAuditCount
     *
     * @return array
     */
    protected function format($config, $withdrawAuditCount)
    {
        // 当前使用版本
        $currentVersions = $config->usedCodeLibrary ?? [];
        // 最近提审记录
        $latestVersions = $config->latestCodeLibrary ?? [];

        // 小程序上线状态集合
        $onlineStatusArray = array_flip(EnterpriseThirdPartAccount::ONLINE_STATUS);
        // 小程序审核状态集合
        $auditStatusArray = array_flip(EnterpriseThirdPartAccount::AUDIT_STATUS);

        // 审核状态
        $auditStatusStr = '';
        // 上线状态
        $onlineStatusStr = $currentVersions->online_status ?? '';
        // 上线方式
        $onlineType = $currentVersions->online_type ?? EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_AUTOMATIC;
        // 审核版本
        $auditVersion = '';
        // 审核时间
        $auditTime = 0;
        
        $timestamp = $currentVersions->created_at ?? 0;
        $currentVersionTime = !empty($currentVersions->created_at) ? $timestamp : 0;

        if ($latestVersions) {
            $onlineStatusStr = $latestVersions->online_status;
            $onlineType = $latestVersions->online_type;
            // 最新提审记录非当前在用版本 && 非撤审状态
            if (EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE === intval($latestVersions->using)
                && EnterpriseMiniProgramCodeLibrary::WITHDRAW !== $latestVersions->audit_status) {
                $auditVersion = $latestVersions->version;
                $auditStatusStr = $latestVersions->audit_status;
                $auditTime = $latestVersions->created_at ?? 0;
            }
        }

        $auditStatus = $auditStatusArray[$auditStatusStr] ?? 0;
        $onlineStatus = $onlineStatusArray[$onlineStatusStr] ?? 0;
        // 获取小程序分阶段发布相关数据
        $grayReleaseData = $this->selGrayReleaseData($config);

        // 审核失败原因
        $reason = '';
        // 审核失败
        if ($auditStatusStr && EnterpriseMiniProgramCodeLibrary::REFUSED == $auditStatusStr) {
            // 查询审核日志
            $auditLog = EnterpriseMiniProgramAuditLog::query()
                ->where('code_library_id', $latestVersions->id)
                ->first();

            $reason = $auditLog->reason ?? '';
        }

        // 当天撤审次数
        $withdrawalCountDay = $withdrawAuditCount[$config->id]['day_count'] ?? 0;
        // 当月撤审次数
        $withdrawalCountMonth = $withdrawAuditCount[$config->id]['month_count'] ?? 0;

        return [
            'account_id'                => $config->id,
            'app_id'                    => $config->app_id,
            'name'                      => $config->name,
            'type'                      => $config->type,
            'current_version'           => $currentVersions->version ?? '',
            'current_version_time'      => $currentVersionTime,
            'gray_release_version'      => $grayReleaseData['gray_release_version'],
            'grp_status'                => $grayReleaseData['grp_status'],
            'grp_gray_percentage'       => $grayReleaseData['grp_gray_percentage'],
            'gray_release_version_time' => $grayReleaseData['gray_release_version_time'],
            'audit_version'             => $auditVersion,
            'audit_status'              => $auditStatus,
            'audit_time'                => $auditTime,
            'reason'                    => $reason,
            'online_type'               => $onlineType,
            'online_status'             => $onlineStatus,
            'withdrawal_count_day'      => $withdrawalCountDay,
            'withdrawal_count_month'    => $withdrawalCountMonth
        ];
    }

    /**
     * 获取小程序分阶段发布相关数据
     *
     * @param $config
     *
     * @return array
     */
    protected function selGrayReleaseData($config)
    {
        // 线上分阶段发布版本
        $grayReleaseVersion = '';
        // 分阶段发布状态：初始状态
        $grpStatus = 0;
        // 灰度的百分比，1到100的整数
        $grpGrayPercentage = 0;
        $grayReleaseVersionTime = 0;

        // 分阶段发布状态：执行完成、被删除
        $releaseStatus = [
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_COMPLETED,
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_DELETED,
        ];

        // 查询小程序分阶段发布数据
        $localGrayRelease = EnterpriseMiniProgramCodeLibrary::query()
            ->where([
                'release_model' => EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY,
                'using'         => EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE
            ])->whereIn('account_id', [$config['id']])
            ->whereNotIn('grp_status', $releaseStatus)
            ->orderByDesc('id')
            ->first();

        // 分阶段发布状态非执行完成或被删除
        if ($localGrayRelease) {
            $grayReleaseVersion = $localGrayRelease['version'];
            $grpStatus = $localGrayRelease['grp_status'];
            $grpGrayPercentage = $localGrayRelease['grp_gray_percentage'];
            $grayReleaseVersionTime = $localGrayRelease['created_at'];
        }

        return [
            'gray_release_version'      => $grayReleaseVersion,
            'grp_status'                => $grpStatus,
            'grp_gray_percentage'       => $grpGrayPercentage,
            'gray_release_version_time' => $grayReleaseVersionTime
        ];
    }
}
