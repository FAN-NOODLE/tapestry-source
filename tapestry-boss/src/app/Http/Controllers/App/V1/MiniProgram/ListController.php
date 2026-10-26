<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:10
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\OpenPlatform\MiniProgramService;

class ListController extends BaseController
{
    protected $miniProgramService;

    public function __construct(MiniProgramService $miniProgramService)
    {
        parent::__construct();
        $this->miniProgramService = $miniProgramService;
    }

    /**
     * 小程序列表
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index()
    {
        // 接收参数
        $params = request()->only('keywords', 'type', 'audit_status', 'page', 'limit');

        // 小程序列表
        $list = $this->miniProgramService->getMiniPrograms($params);
        $miniPrograms = $list['data'] ?? [];

        // 初始化待返数据
        $data = [];
        if (!empty($miniPrograms)) {
            // 数据格式化
            $data = $this->format($miniPrograms);
        }

        return [
            'data' => $data,
            'meta' => [
                'total' => intval($list['total']),
                'size'  => intval($list['per_page']),
                'page'  => intval($list['current_page'])
            ]
        ];
    }

    /**
     * 小程序列表数据格式化
     *
     * @param array $miniPrograms 小程序列表数据
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    protected function format(array $miniPrograms = [])
    {
        // 初始化待返数据
        $data = [];
        // 企业id集合
        $enterpriseIds = array_unique(array_column($miniPrograms, 'enterprise_id'));
        // 企业列表
        $enterprises = Enterprise::query()
            ->where('status', 1)
            ->whereIn('id', $enterpriseIds)
            ->get(['id', 'name'])
            ->toArray();

        // 企业列表
        $enterprises = array_column($enterprises, null, 'id');
        // 小程序上线状态集合
        $onlineStatusArray = array_flip(EnterpriseThirdPartAccount::ONLINE_STATUS);
        // 小程序审核状态集合
        $auditStatusArray = array_flip(EnterpriseThirdPartAccount::AUDIT_STATUS);

        // 授权记录id集合
        $accountIds = array_column($miniPrograms, 'id');
        // 获取小程序分阶段发布数据
        $grayReleaseData = $this->miniProgramService->getGrayReleaseData($accountIds);
        // 小程序分阶段发布数据
        $grayReleaseData = array_column($grayReleaseData, null, 'account_id');

        // 分阶段发布状态：执行完成、被删除
        $releaseStatus = [
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_COMPLETED,
            EnterpriseMiniProgramCodeLibrary::RELEASE_STATUS_DELETED,
        ];

        foreach ($miniPrograms as $item) {
            // 授权记录id
            $accountId = $item['id'];
            // 当前使用版本
            $currentVersions = $item['used_code_library'] ?? [];
            // 最近提审记录
            $latestVersions = $item['latest_code_library'] ?? [];

            // 企业名称
            $enterpriseName = '';
            if (!empty($enterprises[$item['enterprise_id']])) {
                $enterpriseName = $enterprises[$item['enterprise_id']]['name'];
            }

            // 审核状态
            $auditStatusStr = '';
            // 上线状态
            $onlineStatusStr = $currentVersions['online_status'] ?? '';
            // 上线方式
            $onlineType = $currentVersions['online_type'] ?? EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_AUTOMATIC;
            // 审核版本
            $auditVersion = '';

            if ($latestVersions) {
                $onlineStatusStr = $latestVersions['online_status'];
                $onlineType = $latestVersions['online_type'];
                // 最新提审记录非当前在用版本 && 非撤审状态
                if (EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE === intval($latestVersions['using'])
                    && EnterpriseMiniProgramCodeLibrary::WITHDRAW !== $latestVersions['audit_status']) {
                    $auditVersion = $latestVersions['version'];
                    $auditStatusStr = $latestVersions['audit_status'];
                }
            }

            $auditStatus = $auditStatusArray[$auditStatusStr] ?? 0;
            $onlineStatus = $onlineStatusArray[$onlineStatusStr] ?? 0;

            // 线上分阶段发布版本
            $grayReleaseVersion = '';
            // 分阶段发布状态：初始状态
            $grpStatus = 0;
            // 灰度的百分比，1到100的整数
            $grpGrayPercentage = 0;

            // 分阶段发布数据
            $localGrayRelease = $grayReleaseData[$accountId] ?? [];
            if ($localGrayRelease) {
                list($grayReleaseVersion, $grpStatus, $grpGrayPercentage) = $this->updateLocalGrayRelease(
                    $item,
                    $localGrayRelease,
                    $releaseStatus
                );
            }

            // 组装返回数据
            $data[] = [
                'account_id'           => $accountId,
                'app_id'               => $item['app_id'],
                'name'                 => $item['name'],
                'enterprise_id'        => $item['enterprise_id'],
                'enterprise_name'      => $enterpriseName,
                'type'                 => $item['type'], // 小程序类型（1：商城，2：CRS）
                'current_version'      => $currentVersions['version'] ?? '', // 线上全量发布版本
                'gray_release_version' => $grayReleaseVersion, // 线上分阶段发布版本
                'grp_status'           => $grpStatus, // 分阶段发布状态（0:初始状态 1:执行中 2:暂停中 3:执行完毕）
                'grp_gray_percentage'  => $grpGrayPercentage, // 灰度的百分比，1到100的整数
                'audit_version'        => $auditVersion, // 审核版本
                'audit_status'         => $auditStatus, // 审核状态（1：等待审核，2：审核中，3：审核通过，4：审核失败）
                'online_type'          => $onlineType, // 上线方式（1：自动上线，2：手动上线）
                'online_status'        => $onlineStatus, // 上线状态（1：等待上线，2：上线，3：下线）
                'current_tag'          => $currentVersions['release_tag'] ?? '',
                'current_branch'       => $currentVersions['release_branch'] ?? '',
                'latest_tag'           => $latestVersions['release_tag'] ?? '',
                'latest_branch'        => $latestVersions['release_branch'] ?? '',
            ];
        }

        return $data;
    }

    /**
     * 更新本地分阶段发布数据
     *
     * @param array $item
     * @param array $localGrayRelease
     * @param array $releaseStatus
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function updateLocalGrayRelease(array $item, array $localGrayRelease, array $releaseStatus)
    {
        // 线上分阶段发布版本
        $grayReleaseVersion = '';
        // 分阶段发布状态：初始状态
        $grpStatus = 0;
        // 灰度的百分比，1到100的整数
        $grpGrayPercentage = 0;

        // 查询小程序分阶段发布详情
        $grayReleasePlan = $this->miniProgramService->getGrayReleasePlan($item['app_id']);
        $grayReleasePlan = $grayReleasePlan['gray_release_plan'] ?? [];

        if (!empty($grayReleasePlan)) {
            // 分阶段发布状态非执行完成或被删除
            if (!in_array($grayReleasePlan['status'], $releaseStatus)) {
                $grayReleaseVersion = $localGrayRelease['version'];
                $grpStatus = $grayReleasePlan['status'];
                $grpGrayPercentage = $grayReleasePlan['gray_percentage'];
            }

            // 本地与查询的状态不一致，更新本地状态
            if ($grayReleasePlan['status'] != $localGrayRelease['grp_status']) {
                // 待更新数据
                $updateData = ['grp_status' => $grayReleasePlan['status']];

                // 分阶段发布时间+15天
                $grayReleaseEndTime = $grayReleasePlan['create_timestamp'] + 15 * 24 * 60 * 60;
                if ($grayReleaseEndTime >= time()) {
                    // 当前在使用版本改为非当前使用版本
                    EnterpriseMiniProgramCodeLibrary::query()
                        ->where([
                            'account_id' => $item['id'],
                            'using'      => EnterpriseMiniProgramCodeLibrary::USING_STATUS_TRUE
                        ])->update(['using' => EnterpriseMiniProgramCodeLibrary::USING_STATUS_FALSE]);

                    $updateData['using'] = EnterpriseMiniProgramCodeLibrary::USING_STATUS_TRUE;
                    $updateData['release_model'] = EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_ALL;
                }

                $this->miniProgramService->updateGrayReleaseStatus($localGrayRelease['id'], $updateData);
            }
        }

        return [$grayReleaseVersion, $grpStatus, $grpGrayPercentage];
    }
}
