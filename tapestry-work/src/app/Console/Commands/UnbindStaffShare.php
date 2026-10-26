<?php


namespace App\Console\Commands;

use App\Common\Constant;
use App\Common\MessageConstant;
use App\Jobs\Message\ShareLinkMessageAdapter;
use App\Jobs\ShareLinkBatchUpdateCustomerStatisticJob;
use App\Library\Ddc\DdcMessage;
use App\Models\StaffShare;
use App\Services\Events\ShareLinkChangeEventService;
use App\Services\StaffManagerService;
use App\Services\StaffSettingService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UnbindStaffShare extends Command
{
    /**
     * @var string
     */
    protected $signature = 'UnbindStaffShare';

    /**
     * @var string
     */
    protected $description = '解除员工分享绑定关系';

    /**
     * 定时任务
     */
    public function handle()
    {
        Log::info("解除员工分享绑定关系开始", []);

        $enterpriseInfo = StaffShare::query()
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->groupBy('enterprise_id')
            ->get(['enterprise_id'])
            ->toArray();

        // 获取企业id
        $enterpriseIds = array_column($enterpriseInfo, 'enterprise_id');

        if (empty($enterpriseIds)) {
            return true;
        }

        // 批量获取配置
        $settingEnterprise = (new StaffManagerService())->getStaffSetting($enterpriseIds, Constant::LOGIN_IDENTITY_GUIDE);

        if (empty($settingEnterprise)) {
            return true;
        }
        $settingByEnterpriseId = array_column($settingEnterprise, null, 'enterprise_id');

        foreach ($enterpriseIds as $enterpriseId) {
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
            if (empty($enterprise)) {
                Log::info("企业信息获取失败", ['enterprise_id' => $enterpriseId]);
                continue;
            }
            // 获取配置
            $setting = $settingByEnterpriseId[$enterpriseId];

            // 没有配置默认永久有效
            if (empty($setting)) {
                continue;
            }

            // 永久有效
            if ($setting['binding_type'] == StaffShare::BINDING_TYPE_ONE) {
                continue;
            }

            // 解绑时长
            $duration = $setting['duration_time'] * 60 * 60 * 24;

            // 解绑时间节点
            $time = time() - $duration + 60;

            $handleData = StaffShare::query()
                ->where([
                    ['enterprise_id', '=', $enterpriseId],
                    ['binding_type', '=', StaffShare::BINDING_TYPE_YES],
                    ['created_at', '<', $time]
                ])
                ->limit(1000)
                ->get(['id', 'enterprise_id', 'user_id', 'worker_id', 'mem_uid', 'type', 'created_at', 'updated_at'])
                ->toArray();

            if (!empty($handleData)) {
                $ids = array_column($handleData, 'id');

                try {
                    DB::beginTransaction();

                    (new StaffSettingService())->unBindingStaffShareById($ids);

                    DB::commit();

                    Log::info('解除员工分享绑定关系' . $enterpriseId, $ids);

                    $now = time();
                    $shareLinkMessageList = [];
                    $notifyDataList = [];
                    foreach ($handleData as $value) {
                        if ($value['type'] != 10) {
                            $notifyDataList[] = [
                                'enterprise_id' => $value['enterprise_id'],
                                'third_user_id' => $value['user_id'],
                                'curr_worker_id' => $value['worker_id'],
                                'curr_worker_mem_uid' => $value['mem_uid'],
                                'bind_time' => $now
                            ];
                            array_push($shareLinkMessageList, ShareLinkMessageAdapter::deleteMessage($value, $now));
                        }
                    }
                    if (!empty($shareLinkMessageList)) {
                        dispatch(new ShareLinkBatchUpdateCustomerStatisticJob(MessageConstant::SHARE_LINK_BATCH_AUTO_UNBINDING, $shareLinkMessageList));
                    }
                    if (!empty($notifyDataList)) {
                        (new ShareLinkChangeEventService())->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_DELETE, $notifyDataList);
                    }
                } catch (\Exception $e) {
                    DB::rollBack();
                    Log::info("解除员工分享绑定关系异常", [$e->getMessage()]);
                }

            }
        }

        return true;
    }
}
