<?php


namespace App\Console\Commands\Temp;

use App\Common\DateUtils;
use App\Models\DailyMemberCustomerStatistic;
use App\Models\MemberUserRelation;
use App\Models\StoreWorker;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use App\Common\Constant;

class MigrateHistoryMemberCustomerStatistic extends Command
{
    /**
     * @var string
     */
    protected $signature = 'migrate-history-member-customer-statistic';

    /**
     * @var string
     */
    protected $description = '员工历史顾客统计数据迁移';

    /**
     * 使用多进程处理
     */
    public function handle()
    {
        if (!function_exists('pcntl_fork')) {
            exit("pcntl_fork not existing");
        }

        // 初始化子进程数量
        $initSubProcess = 5;
        $handleTime = DateUtils::todayStatTime();

        // 获取今日以前的在职导购id
        $ids = StoreWorker::query()
//            ->where('worker_status', 1)
            ->where('created_at', '<', $handleTime)
            ->orderBy('id')
            ->pluck('id')
            ->toArray();
        if (empty($ids)) {
            return true;
        }

        $workerCount = count($ids);
        // 根据进程数计算数据块大小
        $size = ceil($workerCount / $initSubProcess);
        // 分块
        $workerIds = array_chunk($ids, $size, true);
        // 断开数据库连接，由于子进程会继承主进程的数据库连接，子进程可以通过这个连接读到数据造成数据错乱
        DB::disconnect();

        for ($i = 0; $i < $initSubProcess; $i++) {
            $pid = pcntl_fork();
            if ($pid == -1) {
                echo '子进程创建失败,执行退出.';
                exit;
            } elseif ($pid) {
                echo '父进程PID:' . $pid . PHP_EOL;
                // 防止僵尸进程
                pcntl_wait($status, WNOHANG); // WNOHANG为非阻塞
            } else {
//                echo 'pid:' . posix_getpid() . PHP_EOL;exit;
                // 数据逻辑
                $result = $this->memberDataStatistic($workerIds[$i], $handleTime);
                Log::info('result', $result);
                exit;// 不退出子进程会被子进程再次fork
            }
        }

        // 等待子进程执行结束
        while (pcntl_waitpid(0, $status) != -1) {
            $status = pcntl_wexitstatus($status);
            echo "Child $status completed" . PHP_EOL;
        }
    }

    /**
     * 员工数据处理
     * @param array $workerIds
     * @param int $handleTime
     * @return bool
     */
    public function memberDataStatistic(array $workerIds, $handleTime = 0)
    {
        $result = [
            'success' => [],// 成功
            'fail' => [],// 失败
            'skip' => [],// 已处理跳过
            'non_existent' => [],// 不存在
        ];
        $columns = ['id', 'enterprise_id', 'mem_uid', 'store_id', 'contact_qr_code_id', 'created_at'];
        foreach ($workerIds as $workerId) {
            echo '处理中:' . $workerId . PHP_EOL;
            $storeWorker = StoreWorker::query()
                ->select($columns)
                ->where('id', $workerId)
                ->first();
            if (empty($storeWorker)) { // 不存在的导购
                //Log::info('导购id-' . $workerId . '不存在');
                $result['non_existent'][] = $workerId;
                continue;
            }
            $enterpriseId = $storeWorker->enterprise_id;
            $memUid = $storeWorker->mem_uid;
            // 已处理的mem_uid保存在redis中，避免重复处理
            $successKey = 'member_history_data:success:' . $memUid;

            try {
                $success = Redis::get($successKey);
                if ($success) {
                    $result['skip'][] = $workerId;
                    continue;
                } else {
                    Redis::setex($successKey, 86400, 1);
                }

                // 当日时间
                $sameDay = DateUtils::timeToStatTime(strtotime($storeWorker->created_at));
                // 添加日期距离统计时的时间间隔(天)
                $intervalDays = (int)(($handleTime - $sameDay) / DateUtils::DAY_SECONDS);

                // 员工统计mem_uid
                $cumulativeCustomerNum = 0;// 当前员工累计顾客数
                // 初始化执行的为空表，不需要检测数据是否存在
                /*$sameDayMemberStatistic = DailyMemberCustomerStatistic::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('mem_uid', $memUid)
                    ->where('stat_time', $sameDay)
                    ->first();
                if (!empty($sameDayMemberStatistic)) {
                    $result['skip'][] = $workerId;
                    continue;
                }*/

                $groupData = MemberUserRelation::query()
                    ->select(DB::raw('left(from_unixtime(created_at), 10) as created_date, count(*) as count'))
                    ->where('enterprise_id', $enterpriseId)
                    ->where('mem_uid', $memUid)
                    ->where('delete_type', 0)
                    ->groupBy(DB::raw('left(from_unixtime(created_at), 10)'))
                    ->get()
                    ->toArray();
                $groupData = array_column($groupData, null, 'created_date');

                // 会员数据
                $isMemberData = $this->getMemberData($enterpriseId, $memUid, Constant::IS_MEMBER_STATISTICS);
                $isMemberData = array_column($isMemberData, null, 'created_date');

                // 非会员数据
                $nonMemberData = $this->getMemberData($enterpriseId, $memUid, Constant::NON_MEMBER_STATISTICS);
                $nonMemberData = array_column($nonMemberData, null, 'created_date');

                $cumulativeMemberNum = 0; // 当前员工累计会员数
                $cumulativeNonMemberNum = 0; // 当前员工累计非会员数

                $memberData = [];
                $currentTime = time();
                for ($day = 0; $day < $intervalDays; $day++) {
                    $statTime = $sameDay + ($day * DateUtils::DAY_SECONDS);
                    $date = date('Y-m-d', $statTime);
                    // 查询当日添加的外部联系人
                    $addCount = 0;
                    isset($groupData[$date]) && $addCount = $groupData[$date]['count'];

                    $addMemberCount = 0;
                    isset($isMemberData[$date]) && $addMemberCount = $isMemberData[$date]['count'];

                    $addNonMemberCount = 0;
                    isset($nonMemberData[$date]) && $addNonMemberCount = $nonMemberData[$date]['count'];

                    $cumulativeCustomerNum += $addCount;

                    $cumulativeMemberNum += $addMemberCount;
                    $cumulativeNonMemberNum += $addNonMemberCount;
                    $memberData[] = [
                        'enterprise_id' => $enterpriseId,
                        'mem_uid' => $memUid,
                        'stat_time' => $statTime,
                        'cumulative_customer_num' => $cumulativeCustomerNum,
                        'today_customer_num' => $addCount,
                        'cumulative_member_num' => $cumulativeMemberNum,
                        'today_member_num' => $addMemberCount,
                        'cumulative_non_member_num' => $cumulativeNonMemberNum,
                        'today_non_member_num' => $addNonMemberCount,
                        'created_at' => $currentTime,
                        'updated_at' => $currentTime,
                    ];
                }
                DailyMemberCustomerStatistic::query()
                    ->insert($memberData);

                $result['success'][] = $workerId;
            } catch (\Exception $e) {
//                Log::error($workerId . '处理失败', [$e->getMessage()]);
                Redis::del($successKey);
                $failKey = 'member_history_data:fail:' . $memUid;
                Redis::setex($failKey, 86400, 1);
                $result['fail'][] = $workerId;
            }
        }

        return $result;
    }

    /**
     * @param $enterpriseId
     * @param $memUid
     * @param bool $isMember
     * @return array
     * @Date: 2021/3/29
     * @Time: 15:40
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function getMemberData($enterpriseId, $memUid, $isMember = Constant::IS_MEMBER_STATISTICS)
    {
        return  MemberUserRelation::query()
            ->select(DB::raw('left(from_unixtime(created_at), 10) as created_date, count(*) as count'))
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($isMember), function ($query) use ($isMember) {
                if ($isMember == Constant::IS_MEMBER_STATISTICS) {
                    $query->where('union_code', '<>', '');
                } else if ($isMember == Constant::NON_MEMBER_STATISTICS) {
                    $query->where('union_code', '');
                }
            })
            ->where('mem_uid', $memUid)
            ->where('delete_type', 0)
            ->groupBy(DB::raw('left(from_unixtime(created_at), 10)'))
            ->get()
            ->toArray();
    }
}
