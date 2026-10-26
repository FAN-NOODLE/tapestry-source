<?php


namespace App\Console\Commands;

use App\Common\DateUtils;
use App\Library\Uc\Member;
use App\Models\DailyMemberCustomerStatistic;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;


class InitializeMemberBehaviorData extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'initialize-member-behavior-data {--enterpriseId=}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '初始化导购历史行为数据';


    public function handle()
    {
        Log::info("InitializeMemberBehaviorData.autoExec start", ['time' => time()]);
        $enterpriseId = intval($this->option('enterpriseId'));
        $endStatTime = DateUtils::yesterdayStatTime();
        $syncEnterprise = DailyMemberCustomerStatistic::query()
            ->select('enterprise_id', 'stat_time', 'mem_uid')
            ->when($enterpriseId > 0, function ($query) use($enterpriseId) {
                return $query->where('enterprise_id', $enterpriseId);
            })
            ->where('stat_time', '<=', $endStatTime)
            ->groupBy(['mem_uid'])
            ->orderBy('stat_time', 'ASC')
            ->get();
        $sync = $syncEnterprise->toArray();
        $enterpriseIds = array_column($sync, 'enterprise_id');
        try {
            $enterprises = (new Enterprise())->listNoPage(['enterprise_ids' => array_unique($enterpriseIds)]);
            $enterprises = array_column($enterprises['data'], 'domain', 'id');
        } catch (\Exception $e) {
            Log::error('企业信息查询失败', [
                'enterpriseId' => json_encode($enterpriseIds),
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            return false;
        }
        $member = new Member();
        $time = time();
        foreach ($sync as $v) {
            $memUid = $v['mem_uid'];
            $startStatTime = $v['stat_time'];
            $enterpriseId = $v['enterprise_id'];
            $domain = $enterprises[$enterpriseId] ?? '';
            if (empty($domain)) {
                Log::error('当前企业更新导购信息失败', [
                    'enterpriseId' => $enterpriseId,
                    'domain' => $domain
                ]);
                continue;
            }
            try {
                $startTime = $startStatTime;
                while (true) {
                    if ($startTime > $endStatTime) {
                        break;
                    }
                    $endTime = $startTime + 29 * DateUtils::DAY_SECONDS;
                    if ($endTime > $endStatTime) {
                        $endTime = $endStatTime;
                    }

                    $behavior = $member->getBehaviorData(
                        $domain,
                        [
                            'memUids' => [$memUid],
                            'startTime' => $startTime,
                            'endTime' => $endTime,
                            'dpIds' => []
                        ]
                    );
                    $behavior = is_array($behavior) ? $behavior : [];
                    unset($behavior['requestId']);
                    foreach ($behavior as $b) {
                        $behaviorData = [
                            'new_apply_cnt' => $b['newApplyCnt'] ?? 0,
                            'new_contact_cnt' => $b['newContactCnt'] ?? 0,
                            'chat_cnt' => $b['chatCnt'] ?? 0,
                            'message_cnt' => $b['messageCnt'] ?? 0,
                            'reply_percentage' => ($b['replyPercentage'] ?? 0) * 100,
                            'avg_reply_time' => $b['avgReplyTime'] ?? 0,
                            'negative_feedback_cnt' => $b['negativeFeedbackCnt'] ?? 0,
                            'updated_at' => $time,
                        ];
                        DailyMemberCustomerStatistic::query()->where([
                            'enterprise_id' => $enterpriseId,
                            'stat_time' => $b['statTime'],
                            'mem_uid' => $memUid]
                        )->update($behaviorData);
                    }
                    $startTime = $endTime + DateUtils::DAY_SECONDS;
                }
                $this->info("更新导购信息成功, memUid: {$memUid}, startStatTime: {$startStatTime}");
            } catch (\Exception $e) {
                Log::error('更新当前导购信息失败', [
                    'params' => json_encode($v),
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
                $this->info("更新导购信息失败, memUid:" . $memUid);
                continue;
            }
        }
        Log::info("InitializeMemberBehaviorData.autoExec end", ['time' => time()]);
        return true;
    }
}