<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\MarketingPlan;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskResource;
use App\Models\MarketingPlanVisibleRange;
use App\Models\MemberBannerRelation;
use App\Models\StaffRight;
use App\Models\StoreWorker;
use App\Services\MarketingPlanService;
use App\Services\MemberService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class AutoUpdateMarketingPlanPoster extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:update_marketing_plan_poster';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '营销任务生成/更新海报';

    protected $service;
    protected $enterpriseInfo;

    public function handle()
    {
        Log::info($this->description . '--开始', ['time' => Carbon::now()]);

        //查询有效的计划任务中是小程序海报的素材
        $allTask = MarketingPlanTask::query()
            ->join("marketing_plan_task_resource","marketing_plan_task_resource.task_id","marketing_plan_task.id")
            ->where('marketing_plan_task.end_time',">=",time())
            ->where('marketing_plan_task_resource.type','=',MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER)
            ->where('marketing_plan_task_resource.resource_url','!=','')
            ->where('marketing_plan_task_resource.applet_path','!=','')
            ->where('marketing_plan_task_resource.is_crontab',MarketingPlanTaskResource::IS_CRONTAB_YES)
            ->get(['marketing_plan_task_resource.task_id',"marketing_plan_task.enterprise_id","marketing_plan_task.plan_id"])
            ->toArray();
        if(empty($allTask)){
            Log::info($this->description . '--end', ['time' => Carbon::now()]);
            return true;
        }
        $this->service = new MarketingPlanService();
        $enterpriseIds = array_column($allTask,'enterprise_id');
        //将任务以企业维度组装
        $enterpriseTask = [];
        foreach ($allTask as $k=>$v){
            $enterpriseTask[$v['enterprise_id']][] = $v;
        }
        //以企业维度去处理数据
        foreach ($enterpriseIds as $val){
            $this->enterpriseTask($enterpriseTask[$val]);
        }

        Log::info($this->description . '--end', ['time' => Carbon::now()]);
    }

    /**
     * 已企业维度去执行任务
     * @param $enterpriseTask
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/19 11:28
     */
    public function enterpriseTask($enterpriseTask)
    {
        $enterpriseId = $enterpriseTask[0]['enterprise_id'];
        switchAppHost(['enterprise_id' => $enterpriseId]);
        foreach ($enterpriseTask as $key=>$val){
            //计划是否未执行过定时任务
            //生成二维码及海报
            $this->addOrUpdateRelation($val);
        }
        //如何新增可见范围有数据 处理完清空
        $planIds = array_column($enterpriseTask,'plan_id');
        //如果在执行定时任务期间修改了可见范围,在查询的时候实时生成太阳码和海报
         MarketingPlan::query()
            ->whereIn('id',$planIds)
            ->update(
                [
                    'add_object_ids'=>'',
                    'updated_at'    =>time()
                ]);


        return true;
    }
    /**
     * 更新或新增导购小程序太阳码海报
     * @param $taskInfo
     * @param $status
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 11:39
     */
    public function addOrUpdateRelation($taskInfo)
    {
        $taskID = $taskInfo['task_id'];
        $planID = $taskInfo['plan_id'];
        $resource = MarketingPlanTaskResource::query()
            ->where('task_id',$taskID)
            ->first();
        if(empty($resource)){
            return  true;
        }
        $enterpriseId = $resource->enterprise_id;
        $guideInfo = [];
        //海报背景图是否更新
        $backgroundPosterChange = ($resource->resource_updated == MarketingPlanTaskResource::RESOURCE_UPDATED_YES) ? true : false;
        //小程序路径是否更新
        $pathChange = ($resource->applet_path_update == MarketingPlanTaskResource::APPLET_PATH_YES) ? true : false;

        //新增导购
        //可见范围导购是否新增
        $planInfo = MarketingPlan::query()
            ->where('id',$planID)
            ->first();
        $rangeChange = !empty($planInfo->add_object_ids) ? true : false;
        //新增
        if($rangeChange && !$backgroundPosterChange && !$pathChange){
            $guideInfo =  $this->addPlanRange($planInfo->add_object_ids);
        }else{
            //全部导购
            $guideInfo =  $this->GuideInfoByPlanRange($planID);
        }
        $sunCode = '';
        $qrCodeParams = '';
        $getSunCodeStatus = true;
        foreach ($guideInfo as $guide=>$guides){
            $relation = MemberBannerRelation::query()
                ->where('task_id',$taskID)
                ->where('worker_id',$guides['id'])
                ->first();
            $relationData = empty($relation) ? [] : $relation->toArray();
            //小程序路径更新
            if($pathChange){
                //一种是接口报错 一种是太阳码生成失败(二维码路径错误)     如果太阳码生成失败直接跳出循环 break;
                [$sunCode,$qrCodeParams] =  $this->getSunCode($resource->applet_path,$guides,$enterpriseId,$taskID);
                if(empty($sunCode)){
                    $getSunCodeStatus = false;
                    break;
                }
                //生成海报
                    $poster = $this->service->getPosterImg($sunCode,$resource->resource_url);
                    //创建
                    if(empty($relationData)){
                        $relationData['suncode_url'] = $sunCode;
                        $relationData['suncode_param'] = $qrCodeParams;
                        $relationData['background_url'] = $poster;
                        $relationData['enterprise_id'] = $enterpriseId;
                        $this->addRelationInfo($relationData,$taskID,$guides);
                    }
                    if(!empty($relation)){
                        $updateData = [
                            'suncode_url'   => $sunCode,
                            'suncode_param' => $qrCodeParams,
                            'background_url'=> $poster ?? '',
                            'suncode_mediaid'=> '',
                            'updated_at'    =>time()
                        ];
                        $relation->update($updateData);
                    }
            }
            //小程序路径未更新+背景图更新+任务更新  只生成海报
            //+不是新增的导购
            if(!$pathChange &&  $backgroundPosterChange  && !empty($relation)){
                $poster = $this->service->getPosterImg($relation->suncode_url,$resource->resource_url);
                $updateData = [
                    'background_url'=> $poster ?? '',
                    'suncode_mediaid'=> '',
                    'updated_at'    =>time()
                ];
                $relation->update($updateData);
            }

            //新增的导购
            if(!$pathChange && !$backgroundPosterChange  && empty($relation)){
                //太阳码
                [$sunCode,$qrCodeParams] =  $this->getSunCode($resource->applet_path,$guides,$enterpriseId,$taskID);
                if(empty($sunCode)){
                    $getSunCodeStatus = false;
                    break;
                }
                $poster = $this->service->getPosterImg($sunCode,$resource->resource_url);
                $relationData['suncode_url'] = $sunCode;
                $relationData['suncode_param'] = $qrCodeParams;
                $relationData['background_url'] = $poster;
                $relationData['enterprise_id'] = $enterpriseId;
                $this->addRelationInfo($relationData,$taskID,$guides);
            }
        }
        //如果执行获得太阳阳码没报错
        if($getSunCodeStatus){
            //更新太阳码是否生成状态
            $resource->update([
                'applet_path_update'=> MarketingPlanTaskResource::APPLET_PATH_NO,
                'resource_updated'=> MarketingPlanTaskResource::RESOURCE_UPDATED_NO,
                'is_crontab'    => MarketingPlanTaskResource::IS_CRONTAB_NO,
                'updated_at'    =>time()
            ]);
        }

        return true;
    }

    /**
     * 生成太阳码
     * @param $appletPath
     * @param $guides
     * @param $enterpriseId
     * @param $taskId
     * @return array[]
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/19 14:45
     */
    public function getSunCode($appletPath,$guides,$enterpriseId,$taskId)
    {
            [$sunCode,$qrCodeParams] = $this->service->getSunCode($appletPath,$guides,$enterpriseId,$taskId);
            if(empty($sunCode)){
                $sunCode =  $qrCodeParams = [];
                Log::info("生成任务ID:{$taskId}太阳码失败}");
            }

        return [$sunCode,$qrCodeParams];
    }

    /**
     * 新增导购小程序太阳码海报数据
     * @param $data
     * @param $taskID
     * @param $guides
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 11:21
     */
    public function addRelationInfo($data,$taskID,$guides)
    {
        if(empty($data)){
            return true;
        }
        $insertDate = [
            'enterprise_id' => $data['enterprise_id'],
            'task_id'       => $taskID,
            'mem_uid'       => $guides['mem_uid'],
            'worker_id'     => $guides['id'],
            'suncode_url'   => $data['suncode_url'] ?? '',
            'suncode_param' => $data['suncode_param'] ?? '',
            'background_url'=> $data['background_url'] ?? '',
            'created_at'    =>time(),
            'updated_at'    =>time()
        ];
        MemberBannerRelation::query()
            ->insert($insertDate);
        return true;
    }

    /**
     * 新增的导购数据
     * @param $add_object_ids
     * @param $object_type
     * @param $enterprise_id
     * @return mixed
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 11:52
     */
    public function addPlanRange($add_object_ids)
    {
        //查看更新可见范围的计划 并新增了范围的
        $workerIds = array_values(json_decode($add_object_ids,true));
        $guideInfo = StoreWorker::query()
            ->whereIn('id',$workerIds)
            ->get(['id','mem_uid','store_id'])
            ->toArray();
        return $guideInfo;

    }

    /**
     * 所有导购信息
     * @param $planID
     * @param $enterprise_id
     * @return mixed
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 11:52
     */
    public function GuideInfoByPlanRange($planId)
    {
        //可见范围要以企业维度去查询
        $range = MarketingPlanVisibleRange::query()
            ->where("plan_id",$planId)
            ->get(['worker_id'])
            ->toArray();
        $workerIds = array_column($range,'worker_id');
        $guideInfo = StoreWorker::query()
            ->whereIn('id',$workerIds)
            ->get(['id','mem_uid','store_id'])
            ->toArray();

        return $guideInfo;
    }

}
