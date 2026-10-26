<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Library\Uc\QyWechat;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\MemberExternalBehaviorData;
use App\Models\StoreWorker;
use CrsSdk\Boss\BaseClient;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class SaveUserBehaviorData extends Command
{
    const ALL_DATA_DAYS = 180;

    const MAX_ONCE_DATS = 30;

    const ONCE_SECONDS = self::MAX_ONCE_DATS * 3600 * 24;

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'save-user-behavior-data';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每日导购数据统计（企微获取）';

    /**
     * @throws \Exception
     */
    public function handle($enterpriseId = 0)
    {
        if ($enterpriseId > 0) {
            $enterprises = [(new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data']];
        } else {
            // 查询企业列表
            $enterprises = (new Enterprise())->lists()['data'];
        }
        $yesterday = strtotime('-1 day');
        $now = time();
        $qyWechat = new QyWechat();
        $maxDataSeconds = 3600 * 24 * self::ALL_DATA_DAYS; //180天的秒
        $arr = [];
        foreach ($enterprises as $enterprise) {
            $arr[$enterprise['id']] = 0;
            $memberList = Members::query()->where('enterprise_id', $enterprise['id'])
                ->where('status', Members::MEMBERS_STATUS_ACTIVATED)
                ->get(['id', 'userid'])->toArray();
            $dataOrm = MemberExternalBehaviorData::query()->setConn($enterprise['id']);
            $lastTime = $dataOrm->orderByDesc('stat_time')->value('stat_time');
            // 如果没有数据 或者最新一条数据已经超过180天，则只能获取前180天内数据
            if (empty($lastTime) || (($now - $lastTime) > $maxDataSeconds)) {
                //如果数据为空或者最新数据已经超过180天则获取180天的数据
                $lastTime = $now - $maxDataSeconds;
            }
            $startTime = $lastTime + 24 * 3600;

            while (true) {
                if ($startTime > $yesterday) {
                    break;
                }
                $endTime = $startTime + 24 * (self::MAX_ONCE_DATS - 1) * 3600;
                if ($endTime > $yesterday) $endTime = $yesterday;
                $insertArr = [];
                foreach ($memberList as $member) {
                    $params = [
                        'userid' => [$member['userid']],
                        'start_time' => $startTime,
                        'end_time' => $endTime
                    ];
                    try {
                        $dataList = $qyWechat->getQyResult($params, [
                            'domain' => $enterprise['domain'],
                            'path' => 'cgi-bin/externalcontact/get_user_behavior_data'
                        ]);
                    } catch (\Exception $exception) {
                        Log::info("错误:",[$exception->getCode(),$exception->getMessage()]);
                        if (strstr($exception->getMessage(),'-1:system error')) {
                            sleep(5);
                            $dataList = $qyWechat->getQyResult($params, [
                                'domain' => $enterprise['domain'],
                                'path' => 'cgi-bin/externalcontact/get_user_behavior_data'
                            ]);
                        }
                    }


                    foreach ($dataList['behavior_data'] as $data) {
                        $emptyStatus = true;
                        foreach ($data as $dataKey => $dataVal) {
                            if ($dataKey != 'stat_time' && !empty($dataVal)) $emptyStatus = false;
                        }
                        if ($emptyStatus) continue;
                        $insertArr[] = [
                            'member_id' => $member['id'],
                            'userid' => $member['userid'],
                            'date' => date('Y-m-d', $data['stat_time']),
                            'stat_time' => $data['stat_time'],
                            'new_apply_cnt' => $data['new_apply_cnt'],
                            'new_contact_cnt' => $data['new_contact_cnt'],
                            'chat_cnt' => $data['chat_cnt'],
                            'message_cnt' => $data['message_cnt'],
                            'reply_percentage' => $data['reply_percentage'] ?? -1,
                            'avg_reply_time' => $data['avg_reply_time'] ?? 0,
                            'negative_feedback_cnt' => $data['negative_feedback_cnt'],
                            'created_at' => $now,
                            'updated_at' => $now
                        ];
                    }
                    if (count($insertArr) > 2000) {
                        $dataOrm->insert($insertArr);
                        $arr[$enterprise['id']] += count($insertArr);
                        $insertArr = [];
                    }

                }
                $dataOrm->insert($insertArr);
                $arr[$enterprise['id']] += count($insertArr);
                $startTime = $endTime + 86400;
            }

        }
        Log::info('每日导购数据统计（企微获取）：' . json_encode($arr));
    }


}
