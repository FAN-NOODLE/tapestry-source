<?php
/**
 * Created by PhpStorm.
 * User: Hui Yang
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Common\Constant;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskResource;
use App\Models\MemberBannerRelation;
use App\Services\MarketingPlanService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Carbon;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoUpdateMarketingPlanTaskResource extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:update_marketing_plan_task_resource';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时更新营销计划中素材的media_id';


    public function handle()
    {
        Log::info($this->description . '--开始', ['time' => Carbon::now()]);

        //查询任务结束时间大于现在的任务
        $todayTimestamp = strtotime(date('Y-m-d'));
        $tasks = MarketingPlanTask::query()
            ->where('end_time', '>=', $todayTimestamp)
            ->get()
            ->toArray();
        $taskIdArr = array_column($tasks, 'id');

        //查询计划任务的素材是否需要更新  resource_updated = 1  mediaid_updated_at
        $twoDaysAgoTimeStamp = strtotime("-2 day");

        $processedTaskIdArr = $this->updateOtherResource($taskIdArr, $twoDaysAgoTimeStamp);

        $this->updateBannerResource($processedTaskIdArr);

        Log::info($this->description . '--结束', ['time' => Carbon::now()]);
    }

    /**
     * 处理其它类型，并返回处理了的task_id数组
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/21 10:08
     * @param $taskIdArr
     * @param $limitTime
     * @return array
     * @throws \Exception
     */
    public function updateOtherResource($taskIdArr, $limitTime): array
    {
        //查询需要处理的数据
        $resources = MarketingPlanTaskResource::query()
            ->whereIn('task_id', $taskIdArr)
            ->where('mediaid_updated_at', '<=', $limitTime)
            ->limit(200)
            ->get();
        //查询企业列表
        $enterpriseIds = array_column($resources->toArray(), 'enterprise_id');
        if (empty($enterpriseIds)) {
            return [];
        }
        $enterpriseIds = array_unique($enterpriseIds);
        $enterpriseList = (new Enterprise())->listNoPage(['enterprise_ids' => $enterpriseIds]);
        $enterpriseList = array_column($enterpriseList['data']??[], null, 'id');

        $returnTaskIdArr = [];
        foreach ($resources as $res) {
            if (!isset($enterpriseList[$res->enterprise_id]) || empty($enterpriseList[$res->enterprise_id])) {
                continue;
            }
            switchAppHost(['enterprise_id' => $res->enterprise_id]);
            $qyDomain = $enterpriseList[$res->enterprise_id]['domain'];
            if (!empty($qyDomain)) {
                try {
                    $mediaId = (new MarketingPlanService())->getMediaId(
                        $res->resource_url,
                        $qyDomain,
                        Constant::OPERATE_TYPE_SEND
                    );
                    $res->mediaid = $mediaId;
                    $res->mediaid_updated_at = time();
                    $res->save();
                } catch (\Exception $e) {
                    echo $e->getMessage();
                    continue;
                }
            }
            $returnTaskIdArr[] = $res->task_id;
        }
        return $returnTaskIdArr;
    }

    /**
     * 处理小程序海报
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/21 10:10
     * @param $limitTime
     * @throws \Exception
     */
    public function updateBannerResource($taskIdArr)
    {
        $resources = MemberBannerRelation::query()
            ->where(function ($query) use ($taskIdArr) {
                $query->whereIn('task_id', $taskIdArr)->orWhere('suncode_mediaid', '');
            })
            ->get();
        //查询企业列表
        $enterpriseIds = array_column($resources->toArray(), 'enterprise_id');
        $enterpriseList = (new Enterprise())->listNoPage(['enterprise_ids' => $enterpriseIds]);
        $enterpriseList = array_column($enterpriseList['data']??[], null, 'id');

        foreach ($resources as $res) {
            if (!isset($enterpriseList[$res->enterprise_id]) || empty($enterpriseList[$res->enterprise_id])) {
                continue;
            }
            switchAppHost(['enterprise_id' => $res->enterprise_id]);
            $qyDomain = $enterpriseList[$res->enterprise_id]['domain'];
            if (!empty($qyDomain)) {
                try {
                    $mediaId = (new MarketingPlanService())->getMediaId(
                        $res->background_url,
                        $qyDomain,
                        Constant::OPERATE_TYPE_SEND
                    );
                    $res->suncode_mediaid = $mediaId;
                    $res->updated_at = time();
                    $res->save();
                } catch (\Exception $e) {
                    echo $e->getMessage();
                    continue;
                }
            }
        }
        return true;
    }
}
