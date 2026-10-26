<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 14:23
 */

namespace App\Services;

use App\Common\Constant;
use App\Jobs\SyncNoticeMarketingPlanTask;
use App\Library\ImageTools;
use App\Library\Uc\Department;
use App\Library\Uc\Member;
use App\Models\MarketingPlan;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskProgress;
use App\Models\MarketingPlanTaskResource;
use App\Models\MarketingPlanVisibleRange;
use App\Models\MarketingTaskTagRelation;
use App\Models\MemberBannerRelation;
use App\Models\Scrm\MiniProgram;
use App\Models\Scrm\SourceMiniProgram;
use App\Models\Scrm\WorkExternalTag;
use App\Models\StoreWorker;
use App\Models\SunCodeError;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\AppCode;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;


class MarketingPlanService extends BaseService
{

    protected $size = 10;
    protected $page = 1;
    const PLAN_SN = 'JH';
    /**
     * 素材数据
     *
     * @var
     */
    protected $resourceInfo;
    protected $taskInfo = '';
    protected $imageTools;

    /**
     * 计划列表 获得详情
     * 复制:1 详情：0
     *
     * @var int
     */
    const  DETAIL_COPY = 1;
    const  DETAIL_INFO = 0;

    protected $detail_type = 0;
    /**
     * 当前企业的所有部门
     *
     * @var
     */
    protected $allDepartmentList;
    /**
     * 当前门店的信息
     */
    protected $storeList;
    /**
     *查询导购的字段
     */
    protected $guideSelect = ['id', 'mem_uid', 'store_id'];

    protected $isCrontab = false;

    /**
     * 当前营销计划下在职导购ID集合
     * @var array
     */
    protected $workerIds = [];

    /**
     * 所有导购
     *
     * @var
     */
    protected $guideInfo = [];

    /**
     * 计划详情导购完成任务数据
     * @var array
     */
    protected $completeInfo = [];

    /**
     * 可见范围类型
     *
     * @var
     */
    protected $objectType = 0;
    /**
     * 分享商品类型
     */
    const  SHARE_TYPE = 'isg';

    /**
     * 获取素材media_id
     *
     * @param $im
     * @param $domain
     * @param $optionType
     *
     * @return mixed
     * @throws \Exception
     */
    public function getMediaId($fileName, $domain, $optionType = Constant::OPERATE_TYPE_SEND)
    {
        if ($optionType == Constant::OPERATE_TYPE_DOWNLOAD) {
            ob_start();
            imagepng($fileName);
            $imageData = ob_get_contents();
            ob_end_clean();
            $base64String = base64_encode($imageData);
            //释放空间
            imagedestroy($fileName);

            return $base64String;
//            header('Content-type: image/jpeg');
//            imagejpeg($im,null, 100); // 输出图象到浏览器或文件
        }
        /*if (!is_resource($image)) {
            $image = imagecreatefromstring(file_get_contents($image));
        }
        //输出图片
        $path = base_path('temp');
        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }
        $fileName = $path . '/' . rand(1, 99999) . '.png';
        imagepng($image, $fileName);*/

        $extName  = substr(strrchr($fileName, '.'), 1);
        $extName  = strtolower($extName);
        $fileType = 1;
        //如果是视频
        if (in_array($extName, [
            'mp4',
            'mpeg',
            'avi',
            'mov'
        ])) {
            $fileType = 3;
        }

        $file = @fopen($fileName, 'r');

        // 获取素材id
        $multipart = [
            [
                'name'     => 'file',
                'contents' => $file,
            ],
            [
                'name'     => 'fileType',
                'contents' => $fileType,
            ],
            [
                'name'     => 'filename',
                'contents' => rand(10000, 99999),
            ],
        ];

        $result = (new Member())->imageToMediaId($domain, 'contact', $multipart);

        @fclose($file);

        return is_object($result) ? $result->mediaId : $result['mediaId'];
    }

    /**
     * 创建营销计划
     *
     * @param $params
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 14:15
     */

    public function createPlan($params)
    {
        //创建
        $result = [];
        DB::beginTransaction();
        try {
            //任务表marketing_plan_task
            $this->addTask($params, 0, $params['enterprise_id']);
            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();
            notice(1200003);
            Log::info('创建营销计划失败：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . '.getTrace' . json_encode($exception->getTrace()));
        }

        return $result;
    }

    /**
     * 新增任务
     * @param $params
     * @param $planId
     * @param $enterprise_id
     * @return bool
     */
    public function addTask($params, $planId, $enterprise_id)
    {
        $addResourceInfo = [];
        //客户标签数组
        $tagData = [];
        foreach ($params['task_array'] as $k => $v) {
            //时间处理
//            $start_time        = $this->startOrEndTime($v['start_time'], 1);
//            $end_time          = $this->startOrEndTime($v['end_time'], 2) - 1;
            $start_time        = $v['start_time'] ?? 0;
            $end_time          = $v['end_time'] ?? 0;
            $before_start_time = $v['before_start_time'] ?? 0;
            $before_end_time   = $v['before_end_time'] ?? 0;
            if ($before_start_time) {
                $before_start_time = $start_time - $before_start_time * 60;
            }
            if ($before_end_time) {
                $before_end_time = $end_time - $before_end_time * 60;
            }
            $taskData = [
                'enterprise_id'     => $enterprise_id,
                'task_sn'           => $this->createSN(),
                'plan_id'           => $planId,
                'task_title'        => $v['task_title'],
                'start_time'        => $start_time,
                'end_time'          => $end_time,
                'type'              => $v['type'] ?? 1,
                'before_start_time' => $before_start_time,
                'before_end_time'   => $before_end_time,
                'fixed_time'        => $v['fixed_time'] ?? 0,
                'content'           => $v['content'],
                'created_at'        => time(),
                "updated_at"        => time(),
                'operator_id'       => $params['operator_id'],
                "operator_name"     => $params['operator_name'],
                "object_type"       => $params['object_type'],
            ];

            $taskId = MarketingPlanTask::query()
                ->insertGetId($taskData);

            // 发送提醒
            $this->sendMessage($taskId);

            //可见范围表批量添加 marketing_plan_visible_range
            $this->addPlanVisibleRange($params['object_array'], $params, $taskId);

            //组装关联客户标签入库数据
            if (!empty($v['tag_ids']) && is_array($v['tag_ids'])) {
                foreach ($v['tag_ids'] as $value) {
                    $tagData[] = [
                        'enterprise_id' => $enterprise_id,
                        'task_id'       => $taskId,
                        'tag_id'        => $value,
                        'created_at'    => time(),
                        "updated_at"    => time()
                    ];
                }
            }
            //素材表统一添加
            //拼接素材数据TaskID和企业ID
            foreach ($v['resource'] as $val) {
                $val['enterprise_id'] = $enterprise_id;
                $val['task_id']       = $taskId;
                $addResourceInfo[]    = $val;
            }
        }
        if (!empty($tagData)) {
            MarketingTaskTagRelation::query()->insert($tagData);
        }
        //新增素材表
        $this->addTaskResource($addResourceInfo);

        return true;
    }

    /**
     * 计划编码
     *
     * @return string
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 16:54
     */
    private function createSN()
    {
        $randNumber = rand(0, 9) . rand(0, 9) . rand(0, 9);
        $planSn     = self::PLAN_SN . date('Ymd', time()) . $randNumber;

        return $planSn;
    }

    /**
     * 更新计划
     *
     * @param $params
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:03
     */
    public function updatePlan($params)
    {
        if (empty($params['task_id'])) {
            return true;
        }
        //任务数据 查询任务是否有效
        $planInfo = MarketingPlanTask::query()
            ->with("tags:tag_id,task_id")
            ->where('id', $params['task_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->select(MarketingPlanTask::task_select)
            ->first();
        if (empty($planInfo)) {
            notice(1200002);
        }
        $result = [];
        DB::beginTransaction();
        try {
            //计划表 marketing_plan
            $marketingPlanData = [];
            //$marketingPlanData['plan_title']  = $params['plan_title'];
            $marketingPlanData['object_type'] = $params['object_type'];
            //可见范围表批量添加 marketing_plan_visible_range
            $rangeInfo    = MarketingPlanVisibleRange::query()
                ->where('task_id', $params['task_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->get()
                ->toArray();
            $infoObjectID = array_column($rangeInfo, 'worker_id');
            $objectIDs    = array_column($params['object_array'], 'worker_id');
            //可见范围
            if (!empty($rangeInfo)) {
                //新增的 新数据与旧数据的差集
                $addObjectIDs = array_diff($objectIDs, $infoObjectID);
                if (!empty($addObjectIDs)) {
                    $addObjectData = [];
                    foreach ($params['object_array'] as $val) {
                        if (in_array($val['worker_id'], $addObjectIDs)) {
                            $addObjectData[] = [
                                'worker_id' => $val['worker_id'],
                                'mem_uid'   => $val['mem_uid']
                            ];
                        }
                    }
                    $this->addPlanVisibleRange($addObjectData, $params, $params['task_id']);
                }
                //在计划表中存储新增的部门或门店
                if (!empty($planInfo->add_object_ids)) {
                    $oldIds       = json_decode($planInfo->add_object_ids, true);
                    $addObjectIDs = array_merge($oldIds, $addObjectIDs);
                }
                $marketingPlanData['add_object_ids'] = json_encode($addObjectIDs);

                //移除的 旧数据与新数据的差集
                $removeObjectIDs = array_diff($infoObjectID, $objectIDs);
                if (!empty($removeObjectIDs)) {
                    $this->removePlanVisibleRange($removeObjectIDs, $params['task_id']);
                }
            }
            //新增导购,执行定时任务
            if (!empty($marketingPlanData['add_object_ids'])) {
                $this->isCrontab = true;
            }
            $planInfo->update($marketingPlanData);
            //更新客户标签
            $tagIds = $params['task_array'][0]['tag_ids'] ?? [];
            MarketingTaskTagRelation::query()
                ->where('task_id', $params['task_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->delete();
            if (!empty($tagIds)) {
                $tagData = [];
                foreach ($tagIds as $value) {
                    $tagData[] = [
                        'enterprise_id' => $params['enterprise_id'],
                        'task_id'       => $params['task_id'],
                        'tag_id'        => $value,
                        'created_at'    => time(),
                        "updated_at"    => time()
                    ];
                }
                MarketingTaskTagRelation::query()->insert($tagData);
            }
            //任务表marketing_plan_task
            $this->updateTaskInfo($params['task_array'], $params, $params['task_id']);

            // 发送提醒
            $this->sendMessage($params['task_id']);

            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();
            //dd($exception->getLine());
            notice(1200003);
            Log::info('更新营销计划失败：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . '.getTrace' . json_encode($exception->getTrace()));
        }

        return $result;
    }

    /**
     * 可见范围新增
     *
     * @param $data
     * @param $params
     * @param $taskId
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 16:12
     */
    public function addPlanVisibleRange($data, $params, $taskId)
    {
        if (empty($data)) {
            return true;
        }
        $range = [];
        foreach ($data as $k => $v) {
            $range[] = [
                'enterprise_id' => $params['enterprise_id'],
                'task_id'       => $taskId,
                'worker_id'     => $v['worker_id'],
                'mem_uid'       => $v['mem_uid'],
                'created_at'    => time(),
                'updated_at'    => time()
            ];
        }
        MarketingPlanVisibleRange::query()
            ->insert($range);
        return true;
    }

    /**
     * 移除可见范围
     *
     * @param $data
     * @param $planId
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 16:19
     */
    public function removePlanVisibleRange($infoObjectIDs, $planId)
    {
        MarketingPlanVisibleRange::query()
            ->when(!empty($infoObjectIDs), function ($query) use ($infoObjectIDs) {
                $query->whereIn('worker_id', $infoObjectIDs);
            })
            ->where('task_id', $planId)
            ->delete();

        return true;
    }

    /**
     * 更新任务表
     *
     * @param $taskArray
     * @param $params
     * @param $planId
     *
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:03
     */
    public function updateTaskInfo($taskArray, $params, $planId)
    {
        $taskArray[0]['task_id'] = $planId;
        $taskInfo                = MarketingPlanTask::query()
            ->where('id', $planId)
            ->where('enterprise_id', $params['enterprise_id'])
            ->get()
            ->toArray();
        $taskIDs                 = array_column($taskInfo, 'id');
        $newTaskIDs              = array_column($taskArray, 'task_id');
        //移除的任务ids
        $removeTaskIDs = array_diff($taskIDs, $newTaskIDs);
        if (!empty($removeTaskIDs)) {
            $this->removeTask($removeTaskIDs);
        }
        //新增的任务数据
        //$addTaskInfo    = [];
        $updateTaskData = [];
        foreach ($taskArray as $k => $v) {
            $start_time        = $v['start_time'] ?? 0;
            $end_time          = $v['end_time'] ?? 0;
            $before_start_time = $v['before_start_time'] ?? 0;
            $before_end_time   = $v['before_end_time'] ?? 0;
            if ($before_start_time) {
                $before_start_time = $start_time - $before_start_time * 60;
            }
            if ($before_end_time) {
                $before_end_time = $end_time - $before_end_time * 60;
            }

            //更新的
            $updateTaskData[] = [
                'id'                => $planId,
                'task_title'        => $v['task_title'],
                'start_time'        => $start_time,
                'end_time'          => $end_time,
                'type'              => $v['type'] ?? 1,
                'before_start_time' => $before_start_time,
                'before_end_time'   => $before_end_time,
                'fixed_time'        => $v['fixed_time'] ?? 0,
                'content'           => $v['content'],
                'updated_at'        => time(),
                'operator_id'       => $params['operator_id'] ?? 0,
                "operator_name"     => $params['operator_name'] ?? 0,
                "object_type"       => $params['object_type'] ?? 0,
            ];

            //更新素材
            $this->updateTaskResource($v, $params['enterprise_id']);
        }
        $count = count($updateTaskData);
//        //批量更新任务
//        if (!empty($updateTaskData) && $count > 1) {
//            $this->updateBatch('marketing_plan_task', $updateTaskData);
//        }
        //单条数据更新
        if ($count == 1) {
            $id = $updateTaskData[0]['id'];
            unset($updateTaskData[0]['id']);
            MarketingPlanTask::query()
                ->where('id', $id)
                ->update($updateTaskData[0]);
        }
        //批量新增任务
//        if (!empty($addTaskInfo)) {
//            $this->addTask($addTaskInfo, $planId, $params['enterprise_id']);
//        }

        return true;

    }

    /**
     * 批量更新数据库
     * @param $tableName
     * @param array $multipleData
     * @param string $where
     * @return bool|int
     */
    public function updateBatch($tableName, $multipleData = [], $where = '')
    {
        $updateSql = '';
        try {
            if (empty($multipleData)) {
                return true;
            }
            $firstRow     = current($multipleData);
            $updateColumn = array_keys($firstRow);
            // 默认以id为条件更新，如果没有ID则以第一个字段为条件
            $referenceColumn = isset($firstRow['id']) ? 'id' : current($updateColumn);
            unset($updateColumn[0]);
            // 拼接sql语句
            $updateSql = "UPDATE " . $tableName . " SET ";
            $sets      = [];
            $bindings  = [];
            foreach ($updateColumn as $uColumn) {
                $setSql = "`" . $uColumn . "` = CASE ";
                foreach ($multipleData as $data) {
                    $setSql     .= "WHEN `" . $referenceColumn . "` = ? THEN ? ";
                    $bindings[] = $data[$referenceColumn];
                    $bindings[] = $data[$uColumn];
                }
                $setSql .= "ELSE `" . $uColumn . "` END ";
                $sets[] = $setSql;
            }
            $updateSql .= implode(', ', $sets);
            $updateSql = rtrim($updateSql, ", ") . $where;
            // 传入预处理sql语句和对应绑定数据
            return DB::update($updateSql, $bindings);
        } catch (\Exception $e) {
            Log::info('批量保存失败,表名=' . $tableName . ',数据= ' . json_encode($multipleData) . 'sql=' . $updateSql . ',结果=' . $e->getMessage());
            return false;
        }

    }

    /**
     * 移除任务及任务下的素材
     *
     * @param $taskIDs
     * @param $enterprise_id
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 16:49
     */
    public function removeTask($taskIDs)
    {
        //任务
        MarketingPlanTask::query()
            ->whereIn('id', $taskIDs)
            ->delete();
        //任务素材
        MarketingPlanTaskResource::query()
            ->whereIn('task_id', $taskIDs)
            ->delete();
        //删除太阳码和海报数据 member_banner_relation
        MemberBannerRelation::query()
            ->whereIn('task_id', $taskIDs)
            ->delete();
        return true;
    }

    /**
     * 更新任务素材数据
     *
     * @param $data
     * @param $enterprise_id
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:02
     */
    public function updateTaskResource($data, $enterprise_id)
    {
        //移除旧数据中丢弃的
        $this->removeTaskResource($data, $enterprise_id);
        //新增批量信息素材    //更新的素材数据
        $addResourceInfo = $updateResourceInfo = [];

        foreach ($data['resource'] as $key => $val) {
            $resourceId      = $val['resource_id'] ?? 0;
            $val['sun_code'] = $val['sun_code'] ?? '';
            if (empty($resourceId)) {
                $val['enterprise_id'] = $enterprise_id;
                $val['task_id']       = $data['task_id'];
                $addResourceInfo[]    = $val;
                continue;
            }
            $resourceInfo     = $this->resourceInfo[$val['resource_id']];
            $resource_updated = $applet_path_update = 1;//默认更新
            //素材是否更新
            if (!empty($resourceInfo) && $resourceInfo['resource_url'] == $val['resource_url']) {
                $resource_updated = MarketingPlanTaskResource::RESOURCE_UPDATED_NO;
            }
            //小程序路径是否更新
            if (!empty($resourceInfo) && $resourceInfo['applet_path'] == $val['applet_path']) {
                $applet_path_update = MarketingPlanTaskResource::APPLET_PATH_NO;
            }
            //导购海报表更新数据
            $updateRelation = [];
            //预览海报是否变更
            if (!empty($resourceInfo) && $resourceInfo['applet_path'] != $val['applet_path']) {
                $poster_url                        = '';
                $updateRelation['suncode_url']     = '';
                $updateRelation['suncode_param']   = '';
                $updateRelation['background_url']  = '';
                $updateRelation['suncode_mediaid'] = '';
            }
            if (!empty($resourceInfo) && $resourceInfo['sun_code'] != $val['sun_code']) {
                $poster_url = '';
            }
            if (!empty($resourceInfo) && $resourceInfo['resource_url'] != $val['resource_url']) {
                $poster_url                        = '';
                $mediaid                           = true;
                $updateRelation['background_url']  = '';
                $updateRelation['suncode_mediaid'] = '';
            }

            //可见范围修改
            if ($this->isCrontab) {
                $resourceInfo['is_crontab'] = MarketingPlanTaskResource::IS_CRONTAB_YES;
            }
            //海报背景图修改将微信素材ID设置为空
            if (isset($mediaid)) {
                $resourceInfo['mediaid']            = '';
                $resourceInfo['mediaid_updated_at'] = 0;
            }

            //路径背景修改
            if (isset($poster_url)) {
                $resourceInfo['is_crontab'] = MarketingPlanTaskResource::IS_CRONTAB_YES;
                $resourceInfo['poster_url'] = $poster_url;
            }
            //清空导购海报和太阳码数据
            if (!empty($updateRelation) && $val['type'] == MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER) {
                $this->clearPoster($data['task_id'], $updateRelation);
            }
            //素材表 marketing_plan_task_resource
            $updateResourceInfo[] = [
                'id'                    => $val['resource_id'],
                'type'                  => $val['type'],
                'resource_url'          => $val['resource_url'],
                'applet_path'           => trim($val['applet_path']) ?? '',
                'applet_card_title'     => $val['applet_card_title'] ?? '',
                'applet_card_cover_url' => $val['applet_card_cover_url'] ?? '',
                'resource_updated'      => $resource_updated,
                'applet_path_update'    => $applet_path_update,
                //'sun_code'              => trim($val['sun_code']) ?? '',
                "mediaid"               => $resourceInfo['mediaid'],
                "mediaid_updated_at"    => $resourceInfo['mediaid_updated_at'],
                "poster_url"            => $resourceInfo['poster_url'],
                "is_crontab"            => $resourceInfo['is_crontab'],
                'updated_at'            => time()
            ];
        }
        //批量更新素材
        $count = count($updateResourceInfo);
        //批量更新任务
        if (!empty($updateResourceInfo) && $count > 1) {
            $this->updateBatch('marketing_plan_task_resource', $updateResourceInfo);
        }

        //单条数据更新
        if ($count == 1) {
            $id = $updateResourceInfo[0]['id'];
            unset($updateResourceInfo[0]['id']);
            MarketingPlanTaskResource::query()
                ->where('id', $id)
                ->update($updateResourceInfo[0]);
        }
        //批量新增素材
        $this->addTaskResource($addResourceInfo);
        return true;
    }

    /**
     *  清空导购海报和太阳码数据
     *
     * @param $taskId
     * @param $updateRelation
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/2/3 11:20
     */
    public function clearPoster($taskId, $updateRelation)
    {
        $updateRelation['updated_at'] = time();
        MemberBannerRelation::query()
            ->where('task_id', $taskId)
            ->update($updateRelation);
        return true;
    }

    /**
     *  移除旧数据中丢弃的
     *
     * @param $data
     * @param $enterprise_id
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:01
     */
    public function removeTaskResource($data, $enterprise_id)
    {
        //删除 移除的
        $newIDs = array_column($data['resource'], 'resource_id');
        //统一查询
        $this->getResourceByTaskId($data['task_id'], $enterprise_id);
        $ids       = array_column($this->resourceInfo, 'id');
        $removeIDs = array_diff($ids, $newIDs);
        if (empty($removeIDs)) {
            return true;
        }
        MarketingPlanTaskResource::query()
            ->whereIn('id', $removeIDs)
            ->where('task_id', $data['task_id'])
            ->delete();
        return true;
    }

    /**
     *获得任务的素材
     *
     * @param $taskId
     * @param $enterprise_id
     *
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/19 10:48
     */
    public function getResourceByTaskId($taskId, $enterprise_id)
    {
        $taskAll            = MarketingPlanTaskResource::query()
            ->where('task_id', $taskId)
            ->where('enterprise_id', $enterprise_id)
            ->get()
            ->toArray();
        $taskAll            = array_column($taskAll, null, 'id');
        $this->resourceInfo = $taskAll;
    }

    /**
     * 添加素材表
     *
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 17:01
     */
    public function addTaskResource($data)
    {
        //素材表 marketing_plan_task_resource
        //兼容
        if (empty($data)) {
            return true;
        }
        $resourceData = [];
        foreach ($data as $k => $v) {
            $resourceData[] = [
                'enterprise_id'         => $v['enterprise_id'],
                'task_id'               => $v['task_id'],
                'type'                  => $v['type'],
                'source_id'             => $v['id'] ?? 0,
                'resource_url'          => $v['resource_url'],
                'applet_path'           => isset($v['applet_path']) ? trim($v['applet_path']) : '',
                'applet_card_title'     => $v['applet_card_title'] ?? '',
                'applet_card_cover_url' => $v['applet_card_cover_url'] ?? '',
//                'sun_code'              => $v['sun_code'] ?? '',
                'created_at'            => time(),
                'updated_at'            => time()
            ];
        }

        MarketingPlanTaskResource::query()
            ->insert($resourceData);

        return true;
    }

    /**
     * 删除计划
     *
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:23
     */
    public function delete($params)
    {
        $planInfo = MarketingPlanTask::query()
            ->where('id', $params['task_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->get()
            ->toArray();
        $result   = [];
        //删除任务和素材
        DB::beginTransaction();
        try {
            if (!empty($planInfo)) {
                $taskIds = array_column($planInfo, 'id');
                $this->removeTask($taskIds);
            }
            //删除可见范围
            MarketingPlanVisibleRange::query()
                ->where('task_id', $params['task_id'])
                ->delete();
            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();
            $result = [$exception->getMessage()];
            Log::info('删除营销计划失败：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . '.getTrace' . json_encode($exception->getTrace()));
        }

        return $result;
    }

    public function startOrEndTime($time, $type = 1)
    {
        //开始时间
        if ($type == 1) {
            return strtotime(date('Y-m-d', $time));
        }
        //结束时间
        if ($type == 2) {
            return strtotime(date('Y-m-d', $time)) + 86400;
        }
    }

    /**
     * 获得计划列表
     *
     * @param $params
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 14:40
     */
    public function list($params)
    {
        $this->page = $params['page'] ?? 1;
        $this->size = $params['size'] ?? 10;
        $time       = time();
        //将时间处理成开始或结束
        if (!empty($params['create_time']['create_start'])) {
            $params['create_time']['create_start'] = $this->startOrEndTime($params['create_time']['create_start'], 1);
        }
        if (!empty($params['create_time']['create_end'])) {
            $params['create_time']['create_end'] = $this->startOrEndTime($params['create_time']['create_end'], 2);
        }
        if (!empty($params['edit_time']['create_start'])) {
            $params['edit_time']['create_start'] = $this->startOrEndTime($params['edit_time']['create_start'], 1);
        }
        if (!empty($params['edit_time']['edit_end'])) {
            $params['edit_time']['edit_end'] = $this->startOrEndTime($params['edit_time']['edit_end'], 2);
        }

        $list = MarketingPlanTask::query()
            ->with(['range' => function ($query) {
            }])
            ->with(['resource' => function ($query) {
            }])
            ->where('marketing_plan_task.enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['task_title']), function ($query) use ($params) {
                $query->where('task_title', 'like', '%' . trim($params['task_title']) . '%');
            })
            // 任务类型
            ->when(!empty($params['type']), function ($query) use ($params) {
                $query->where('type', $params['type']);
            })
            // 任务状态
            ->when(!empty($params['status']), function ($query) use ($params, $time) {
                if ($params['status'] == MarketingPlanTask::TASK_NOT_STARTED) {
                    return $query->where(
                        'start_time', ">", $time
                    );
                } else if ($params['status'] == MarketingPlanTask::TASK_IN_PROGRESS) {
                    return $query->where([
                        ['start_time', "<=", $time],
                        ['end_time', ">=", $time]
                    ]);
                } else if ($params['status'] == MarketingPlanTask::TASK_COMPLETE) {
                    return $query->where(
                        'end_time', "<", $time
                    );
                }
            })
            //任务编号
            ->when(!empty($params['task_sn']), function ($query) use ($params) {
                $query->where('task_sn', 'like', '%' . trim($params['task_sn']) . '%');
            })
            //创建时间
            ->when(!empty($params['create_time']['create_start']), function ($query) use ($params) {
                $query->where('created_at', ">=", $params['create_time']['create_start']);
            })
            ->when(!empty($params['create_time']['create_end']), function ($query) use ($params) {
                $query->where('created_at', "<", $params['create_time']['create_end']);
            })
            //编辑时间
            ->when(!empty($params['edit_time']['edit_start']), function ($query) use ($params) {
                $query->where('updated_at', ">=", $params['edit_time']['edit_start']);
            })
            ->when(!empty($params['edit_time']['edit_end']), function ($query) use ($params) {
                $query->where('updated_at', "<", $params['edit_time']['edit_end']);
            })
            //创建人
            ->when($params['operator_name'] != '', function ($query) use ($params) {
                $query->where('operator_name', 'like', '%' . trim($params['operator_name']) . '%');
            })
            ->orderBy('id', 'desc')
            ->paginate($this->size, ['*'], '', $this->page)
            ->toArray();

        if (empty($list['data'])) {
            $list['list'] = [];
            $list['page'] = $this->page;
            $list['size'] = $this->size;
            return $list;
        }
        $workerAllIds = [];
        foreach ($list['data'] as $value) {
            $workerIds    = array_column($value['range'], 'worker_id');
            $workerAllIds = array_merge($workerAllIds, $workerIds);
        }
        //获得所有导购数据
        if (!empty($workerAllIds)) {
            $workerAllIds = array_unique($workerAllIds);
            $this->getGuideInfoByIds($workerAllIds);
        }
        $list = $this->formDataList($list);

        return $list;
    }


    /**
     * 详情
     *
     * @param $params
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 15:18
     */
    public function detail($params)
    {
        $this->detail_type = $params['type'];
        switchAppHost(['enterprise_id' => $params['enterprise_id']]);
        //范围
        $rangeInfo = MarketingPlanVisibleRange::query()
            //->where('plan_id', $params['plan_id'])
            ->where('task_id', $params['task_id'])
            ->select(MarketingPlanVisibleRange::range_select)
            ->get()
            ->toArray();
        $workerIds = array_column($rangeInfo, 'worker_id');
        $this->getGuideInfoByIds($workerIds);

        //导购数据按照入库先后排序
        $guideInfo = $this->getObjectNameByIds($workerIds);
        //任务
        $taskInfo = MarketingPlanTask::query()
            ->with("tags:tag_id,task_id")
            ->where('id', $params['task_id'])
            ->select(MarketingPlanTask::task_select)
            ->get()
            ->toArray();

        $taskIds = array_column($taskInfo, 'task_id');
        $tagIds  = array_column($taskInfo[0]['tags'] ?? [], 'tag_id');
        //客户标签
        $tagInfo = WorkExternalTag::query()
            ->whereIn('id', $tagIds)
            ->get()
            ->toArray();
        //素材数据
        $resourceInfo = MarketingPlanTaskResource::query()
            ->whereIn('task_id', $taskIds)
            ->select(MarketingPlanTaskResource::resource_select)
            ->get()
            ->toArray();
        //小程序素材信息
        $miniSourceIds     = collect($resourceInfo)
            ->where('type', MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER)
            ->pluck('source_id')
            ->all();
        $miniTitles        = SourceMiniProgram::query()->whereIn('id', $miniSourceIds)->get();
        $miniIdToTitle     = $miniTitles->pluck('title', 'id')->all();
        $miniIdToProgramId = $miniTitles->pluck('mini_program_id', 'id')->all();
        $miniList          = MiniProgram::query()
            ->whereIn('id', $miniTitles->pluck('mini_program_id')->all())
            ->pluck('app_name', 'id')
            ->all();

        $resourceInfoNew = [];
        foreach ($resourceInfo as $rekey => $resource) {
            if ($resource['type'] == MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER) {
                $resource['title']    = $miniIdToTitle[$resource['source_id']] ?? '';
                $resource['app_name'] = $miniList[$miniIdToProgramId[$resource['source_id']]] ?? '';
            }
            //生成后台预览海报
            if ($resource['type'] == MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER &&
                empty($resource['poster_url'])
            ) {
                [$posterUrl, $sunCode] = $this->createPreviewPoster($resource);
                $resource['poster_url'] = $posterUrl;
                $resource['sun_code']   = $sunCode;
            }
            if (!empty($this->detail_type)) {
                unset($resource['resource_id']);
            }

            $resourceInfoNew[$resource['task_id']][] = $resource;
        }
        //导购任务完成导购离职
        //导购任务完成数据
        $this->getCompleteInfoByTaskIds($taskIds);
        //任务数据
        foreach ($taskInfo as $key => $val) {
            if (empty($this->detail_type)) {
                //复制不展示任务时间
                $val['start_time'];
                $val['end_time'];
                $val['before_start_time'];
                $val['before_end_time'];
                $val['fixed_time'];
                //合伙人任务完成情况
                $effectiveGuideIDs      = $this->calculationTaskCompletion($val['task_id'], $this->workerIds);
                $val['all_amount']      = count($this->workerIds);
                $val['complete_amount'] = count($effectiveGuideIDs);
                $val['resource']        = $resourceInfoNew[$val['task_id']] ?? [];
            }
            $taskInfo[$key]                 = $val;
            $taskInfo[$key]['tags']         = $tagInfo;
            $taskInfo[$key]['object_array'] = $guideInfo ?? [];
            $taskInfo[$key]['sum']          = $taskInfo[0]['all_amount'] ?? 0;
            if ($taskInfo[$key]['before_start_time']) {
                $taskInfo[$key]['before_start_time'] = ($taskInfo[$key]['start_time'] - $taskInfo[$key]['before_start_time']) / 60;
            }
            if ($taskInfo[$key]['before_end_time']) {
                $taskInfo[$key]['before_end_time'] = ($taskInfo[$key]['end_time'] - $taskInfo[$key]['before_end_time']) / 60;
            }
        }

        return $taskInfo[0] ?? [];
    }

    /**
     * 计划任务导购的完成数据
     * @param $taskIds
     * @param $workerId
     * @return array
     */
    public function getCompleteInfoByTaskIds($taskIds, $workerId = '')
    {
        //导购任务完成数据
        $completeInfo = MarketingPlanTaskProgress::query()
            ->whereIn('task_id', $taskIds)
            ->orderBy('created_at', 'DESC')
            ->when(!empty($workerId), function ($query) use ($workerId) {
                return $query->where('worker_id', $workerId);
            })
            ->get()
            ->toArray();
        foreach ($completeInfo as $k => $v) {
            $this->completeInfo[$v['task_id']][] = $v;
        }
        return array_keys($this->completeInfo);
    }

    /**
     * 任务完成数据
     *
     * @param $planID
     * @param $taskID
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 17:28
     */
    public function calculationTaskCompletion($taskId, $workerIds)
    {
        //完成的导购（完成中有已移除的）
        //marketing_plan_task_progress
        $completeInfo     = $this->completeInfo[$taskId] ?? [];
        $completeGuideIDs = array_column($completeInfo, 'worker_id');
        //全部导购 $workerIds
        //已完成:求交集就是已完成的数量
        $effectiveGuideIDs = array_intersect($workerIds, $completeGuideIDs);

        return $effectiveGuideIDs;
    }


    /**
     * @param $list
     *
     * @return mixed
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 16:25
     */
    public function formDataList($list)
    {
        $dataNew = [];

        foreach ($list['data'] as $key => $val) {
            $workerIds = array_column($val['range'], 'worker_id');
            $rangeData = $this->getObjectNameByIds($workerIds);
            //素材
            $resourceType = [];
            foreach ($val['resource'] as $resource) {
                $resourceType[] = $resource['type'];
            }
            $resourceType = array_values(array_unique($resourceType));

            $status = MarketingPlanTask::TASK_NOT_STARTED;
            $time   = time();
            if ($val['start_time'] <= $time && $val['end_time'] >= $time) {
                $status = MarketingPlanTask::TASK_IN_PROGRESS;
            } else if ($val['end_time'] < $time) {
                $status = MarketingPlanTask::TASK_COMPLETE;
            }

            $dataNew[] = [
                "task_id"       => $val['id'],
                "task_sn"       => $val['task_sn'],
                "task_title"    => $val['task_title'],
                "type"          => $val['type'] ?? MarketingPlanTask::TYPE_TRAINING_TASK,
                "status"        => $status,
                "object_type"   => $val['object_type'] ?? 1,
                "object_array"  => $rangeData,
                "plan_resource" => $resourceType,
                "created_at"    => $val['created_at'],
                "updated_at"    => $val['updated_at'],
                "operator_name" => $val['operator_name']

            ];
        }
        unset($list['data']);
        $list['list'] = $dataNew;
        $list['page'] = $this->page;
        $list['size'] = $this->size;
        return $list;
    }

    /**
     * 导购信息和小程序路径获得太阳码
     *
     * @param $applet_path
     * @param $guideInfo
     * @param $enterprise_id
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 14:47
     */
    public function getSunCode($applet_path, $guideInfo, $enterprise_id, $taskId = '')
    {
        //小程序路径和参数
        [$path, $params] = $this->convertUrlQuery($applet_path);
        //参数如果是加密先解密在加密
        if (isset($params['scene']) && !empty($params['scene'])) {
            $decryptData = decryptShortLink($params['scene']);
            $params      = [];
            foreach ($decryptData['args'] ?? [] as $key => $value) {
                $params[$value['name']] = $value['value'];
            }
        }
        //导购参数
        $qrCodeParams = [
            'store_id'  => $guideInfo['store_id'] ?? 0,
            'source'    => self::SHARE_TYPE,
            'share_uid' => $guideInfo['mem_uid'] ?? '',
        ];
        //二维码参数
        $qrCodeParams = array_merge($qrCodeParams, $params);
        try {
            $sunCode = (new AppCode())->create(['enterprise_id' => $enterprise_id, 'params' => $qrCodeParams, 'path' => $path])['data']['code_url'] ?? '';
            if (empty($sunCode)) {
                $data = [
                    'error'          => 'shop接口返回数据为空或小程序路径有误',
                    'sun_code_param' => json_encode($qrCodeParams),
                    'task_id'        => $taskId ?? '',
                    'enterprise_id'  => $enterprise_id,
                    "created_at"     => time(),
                    "updated_at"     => time()
                ];
                SunCodeError::query()
                    ->insert($data);
            }
        } catch (\Exception $exception) {
            $sunCode = [];
            //记录sql
            $qrCodeParams['path'] = $path;
            $data                 = [
                'error'          => json_encode($exception->getMessage()),
                'sun_code_param' => json_encode($qrCodeParams),
                'task_id'        => $taskId ?? '',
                'enterprise_id'  => $enterprise_id,
                "created_at"     => time(),
                "updated_at"     => time()
            ];
            SunCodeError::query()
                ->insert($data);
            Log::info('计划任务太阳码生成：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . '.getTrace' . json_encode($exception->getTrace()));
        }
        Log::info('计划任务太阳码生成', [$sunCode]);
        //小程序太阳码参数包括小程序路径
        $qrCodeParams['path'] = $path;

        return [$sunCode, json_encode($qrCodeParams)];
    }

    /**
     * 小程序路径解参
     *
     * @param $applet_path
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 15:03
     */
    public function convertUrlQuery($applet_path)
    {
        $path   = '';
        $params = [];
        if (empty($applet_path)) {
            return [$path, $params];
        }

        $tmp   = explode("?", $applet_path);
        $path  = $tmp[0];
        $query = $tmp[1] ?? '';
        if (!empty($query)) {
            parse_str($query, $params);
        }
        return [$path, $params];
    }

    /**
     * 获得小程序海报
     *
     * @param $sunCode
     * @param $backgroundImg
     *
     * @return mixed|string
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 16:24
     */
    public function getPosterImg($sunCode, $backgroundImg)
    {
        $this->imageTools = new ImageTools();
        //背景图本地图片
        $backgroundImg = $this->imageTools->getImage($backgroundImg);

        $localPoster = '';
        if (pathinfo($backgroundImg)['extension'] == 'jpeg') {
            $localPoster = $this->posterByImgJpg($sunCode, $backgroundImg);
        }
        //png 合成有黑边需要单独处理
        if (pathinfo($backgroundImg)['extension'] == 'png') {
            $localPoster = $this->posterByImgPng($sunCode, $backgroundImg);
        }
        //合伙海报
        unlink($backgroundImg);
        $url = $this->uploadImgCloud($localPoster);
        return $url;
    }

    /**
     * 合成海报/两张图片合成
     *
     * @param $sunCode
     * @param $backgroundImg
     *
     * @return string
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 16:21
     */
    public function posterByImgJpg($sunCode, $backgroundImg)
    {
        // 海报背景图片
        $bg_file = base_path('public/sourceImages/background@1x.png');
        $bg_file = imagecreatefromstring(file_get_contents($bg_file));
        //太阳码资源
        $sunCode = $this->imageTools->getImageSource($sunCode)[0];
        //字体文件
        // $font_file = base_path('public/fonts/msyh.ttf');
        $im = imagecreatefromjpeg($backgroundImg);
        //太阳码背景图
        $text_w = imagesx($bg_file);
        $text_h = imagesy($bg_file);
        imagecopyresampled($im, $bg_file, 500, 926, 0, 0, $text_w, $text_h, $text_w, $text_h);
        //太阳码
        $sunCode_w = imagesx($sunCode);
        $sunCode_h = imagesy($sunCode);
        imagecopyresampled($im, $sunCode, 536, 946, 0, 0, 128, 128, $sunCode_w, $sunCode_h);

        //输出图片
        $path = base_path('temp');
        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }
        $fileName = $path . '/' . rand(1, 99999) . '.png';
        imagepng($im, $fileName);
        return $fileName;
    }

    /**
     * PNG 多图片合成黑边问题
     *
     * @param $background
     * @param $sunCode
     *
     * @return string
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/28 10:51
     */
    public function posterByImgPng($sunCode, $background)
    {
        // 海报背景图片
        $bg_file = base_path('public/sourceImages/background@1x.png');
        $bg_file = imagecreatefromstring(file_get_contents($bg_file));

        $background = imagecreatefromstring(file_get_contents($background));
        $sunCode    = imagecreatefromstring(file_get_contents($sunCode));
        //太阳码背景图
        $text_w = imagesx($bg_file);
        $text_h = imagesy($bg_file);
        $cut    = imagecreatetruecolor($text_w, $text_h);
        imagecopy($cut, $background, 0, 0, 500, 926, $text_w, $text_h);
        $opacity = 100;
        imagecopy($cut, $bg_file, 0, 0, 0, 0, $text_w, $text_h);
        imagecopymerge($background, $cut, 500, 926, 0, 0, $text_w, $text_h, $opacity);

        $sun_w = imagesx($sunCode);
        $sun_h = imagesy($sunCode);
        imagecopyresampled($background, $sunCode, 536, 946, 0, 0, 128, 128, $sun_w, $sun_h);
        //输出图片
        $path = base_path('temp');
        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }
        $fileName = $path . '/' . rand(1, 99999) . '.png';
        imagepng($background, $fileName);
        return $fileName;
    }

    /**
     * 上传图片到云删除本地图片
     *
     * @param $urlLocalPath
     *
     * @return mixed
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 16:18
     */
    public function uploadImgCloud($urlLocalPath)
    {
        if (empty($urlLocalPath)) {
            return '';
        }
        $result = CommonService::uploadImage($urlLocalPath);
        //        $result = (new AzureBlobService())->uploadByFilePath($urlLocalPath);
        // 保存生成的海报
        if ($result) {
            unlink($urlLocalPath);
        }
        return $result;
    }


    /**
     * 获得计划可见范围的导购
     *
     * @param $object_type
     * @param $infoObjectIDs
     * @param $enterprise_id
     *
     * @return array
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/14 15:51
     */
    public function GuideInfoByPlanRange($object_type, $infoObjectIDs, $enterprise_id, $params)
    {
        if (empty($object_type)) {
            return [];
        }
        $guideInfo = [];
        //区域维度
        if ($object_type == MarketingPlanTask::OBJECT_TYPE_STORE) {
            $guideInfo = $this->getGuideByStoreIDs($infoObjectIDs, $enterprise_id);
        }
        //部门维度
        if ($object_type == MarketingPlanTask::OBJECT_TYPE_DEPARTMENT) {
            $guideInfo = $this->getGuideByDepartment($infoObjectIDs, $enterprise_id, $params);
        }

        return $guideInfo;
    }

    /**
     * 部门下导购
     *
     * @param $infoObjectIDs
     * @param $enterprise_id
     *
     * @return array
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/14 15:50
     */
    public function getGuideByDepartment($infoObjectIDs, $enterprise_id, $enterpriseInfo)
    {
        // 调UC获取部门下的人员列表
        $qyDomain           = $enterpriseInfo['qy_domain'] ?? '';
        $appIdentifier      = $enterpriseInfo['app_identifier'] ?? '';
        $params['dpIdList'] = $infoObjectIDs;
        $params['pageNum']  = 1;
        $params['pageSize'] = 9999;
        //查询子目录
        //$params['departmentChildrenFlag'] = 1;
        $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);
        // 获取人员的memUid
        $memUids = array_column($result['list'], 'memUid');
        $memUids = array_unique($memUids);
        // 查询人员对应的导购信息
        $guideInfo = StoreWorker::query()
            ->whereIn('mem_uid', $memUids)
            ->where('enterprise_id', $enterprise_id)
            ->where('worker_status', StoreWorker::IN_SERVICE)
            ->select($this->guideSelect)
            ->get()
            ->toArray();
        if ($guideInfo) {
            // 获取人员和部门对应信息
            $result = array_column($result['list'], 'dpName', 'memUid');
            // 加入部门信息
            foreach ($guideInfo as &$value) {
                $value['store_name'] = isset($result[$value['mem_uid']][0]['dpName']) ? $result[$value['mem_uid']][0]['dpName'] : '';
            }
        }
        return $guideInfo;
    }

    /**
     * 依据门店获得导购信息
     *
     * @param $infoObjectIDs
     * @param $enterprise_id
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/14 14:43
     */
    public function getGuideByStoreIDs($infoObjectIDs, $enterprise_id)
    {
        $guideInfo = StoreWorker::query()
            ->whereIn('store_id', $infoObjectIDs)
            ->where('enterprise_id', $enterprise_id)
            ->where('worker_status', StoreWorker::IN_SERVICE)
            ->select($this->guideSelect)
            ->get()
            ->toArray();
        return $guideInfo;
    }

    public function getMarketingPlanTaskCalendarList($params)
    {
        //将结束时间定义为当天的最后一秒
        $params['end_time'] = date('Y-m-d 23:59:59', $params['end_time']);
        $params['end_time'] = strtotime($params['end_time']);

        $result = MarketingPlanTask::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('marketing_plan_task.plan_id', $params['plan_ids'])
            ->where(function ($query) use ($params) {
                $betweenDateArr = [
                    $params['start_time'],
                    $params['end_time']
                ];
                $query->whereBetween('start_time', $betweenDateArr)
                    ->orWhereBetween('end_time', $betweenDateArr)
                    ->orWhere(function ($query) use ($params) {
                        $query->where('start_time', '<=', $params['start_time'])
                            ->where('end_time', '>=', $params['end_time']);
                    });
            })
            ->orderBy('id', 'DESC')
            ->select([
                'id',
                'plan_id',
                'task_title',
                'start_time',
                'end_time',
                'content'
            ])
            ->get()
            ->toArray();

        return $result;
    }

    /**
     * 获取任务列表
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/25 19:49
     *
     * @param $params
     *
     * @return array
     */
    public function getMarketingPlanTaskList($params)
    {
        //将结束时间定义为当天的最后一秒
        //$params['end_time'] = date('Y-m-d 23:59:59', $params['end_time']);
        //$params['end_time'] = strtotime($params['end_time']);

        $fields = [
            'marketing_plan_task.id',
            'marketing_plan_task.plan_id',
            'marketing_plan_task.task_title',
            'marketing_plan_task.start_time',
            'marketing_plan_task.end_time',
            'marketing_plan_task.content',
            'marketing_plan_task.type as task_type'
        ];

        if (in_array($params['is_finished'], [0, MarketingPlanTaskProgress::TASK_ALREADY_FINISHED])) {
            $fields = array_merge($fields, [DB::raw('if(marketing_plan_task_progress.id != "", 2, 1) as is_finished')]);
        } else {
            $fields = array_merge($fields, [DB::raw('1 as is_finished')]);
        }

        $pageList = MarketingPlanTask::query()
            ->with([
                'resource' => function ($query) use ($params) {
                    $query->select([
                        'id',
                        'task_id',
                        'type',
                        'mediaid',
                        'resource_url',
                        'applet_path',
                        'applet_card_title',
                        'applet_card_cover_url'
                    ])->orderBy('id', 'asc');
                }])
            //过滤掉没有生成海报的数据，没有获取到生成media_id的数据
//            ->whereNotExists(function ($query) {
//                $rawSql = 'marketing_plan_task.id=marketing_plan_task_resource.task_id';
//                $rawSql .= ' and marketing_plan_task_resource.mediaid=""';
//                $rawSql .= ' and marketing_plan_task_resource.deleted_at=0';
//                $query->select(DB::raw(1))
//                    ->from('marketing_plan_task_resource')
//                    ->whereRaw($rawSql);
//            })
            ->when(in_array($params['is_finished'], [0, MarketingPlanTaskProgress::TASK_ALREADY_FINISHED]), function ($query) use ($params) {
                $selectType = $params['is_finished'] == MarketingPlanTaskProgress::TASK_ALREADY_FINISHED ? 'inner' : 'left';

                $query->join('marketing_plan_task_progress', function ($join) use ($params) {
                    $join->on('marketing_plan_task.id', '=', 'marketing_plan_task_progress.task_id')
                        ->where('marketing_plan_task_progress.worker_id', '=', $params['worker_id']);
                }, null, null, $selectType);
            })
            ->where('marketing_plan_task.enterprise_id', $params['enterprise_id'])
            ->whereIn('marketing_plan_task.id', $params['task_ids'])
            ->when(!empty($params['task_type']), function ($query) use ($params) {
                $query->where('marketing_plan_task.type', $params['task_type']);
            })
            ->when(!empty($params['task_title']), function ($query) use ($params) {
                $query->where('marketing_plan_task.task_title', 'like', '%' . $params['task_title'] . '%');
            })
            //->whereIn('marketing_plan_task.plan_id', $params['plan_ids'])

            ->when(isset($params['is_index']), function ($query) use ($params) {
                if ($params['is_index']) {
                    $query->where('marketing_plan_task.start_time', '<=', strtotime(date('Y-m-d 23:59:59', $params['current_date'])))
                        ->where('marketing_plan_task.end_time', '>=', $params['current_date']);
                } else {
                    if (!empty($params['start_time']) && !empty($params['end_time'])) {
                        $query->where([
                            ['marketing_plan_task.start_time', "<=", $params['end_time']],
                            ['marketing_plan_task.end_time', ">=", $params['start_time']]
                        ]);
                    }
                }
            })
            ->orderBy('is_finished', 'ASC')
            ->when(in_array($params['is_finished'], [0, MarketingPlanTaskProgress::TASK_ALREADY_FINISHED]), function ($query) use ($params) {
                $query->orderBy('marketing_plan_task_progress.created_at', 'DESC');
            })
            ->orderBy('marketing_plan_task.end_time', 'ASC')
            ->paginate($params['size'], $fields)
            ->toArray();
        $total    = $pageList['total'];
        $result   = $pageList['data'];

        if (empty($result)) {
            return [
                $result,
                $total
            ];
        }

        $taskIdArray = array_column($result, 'id');
        $banners     = MemberBannerRelation::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('task_id', $taskIdArray)
            ->where('worker_id', $params['worker_id'])
            ->get([
                'suncode_mediaid',
                'background_url',
                'task_id'
            ])->toArray();
        $banners     = array_column($banners, null, 'task_id');

        //这个appid是发小程序卡片要用的
//        $shopDetail = (new MiniProgram())->detail([
//            'enterprise_id' => $params['enterprise_id']
//        ]);
//        $appId      = $shopDetail['data']['app_id'] ?? '';
        $appId = '';

        $traceParams['mem_uid']   = $params['mem_uid'] ?? '';
        $traceParams['share_uid'] = $params['mem_uid'] ?? '';
        $traceParams['source']    = self::SHARE_TYPE;
        $traceParams['worker_id'] = $params['worker_id'] ?? 0;
        $traceParams['store_id']  = $params['store_id'];

        //处理格式
        foreach ($result as &$val) {
            //资源的类型只会有一种，将type加在外面，方便前端使用
            $val['type'] = 0;
            if (!empty($val['resource'])) {
                $firstItem   = current($val['resource']);
                $val['type'] = $firstItem['type'];

                if ($val['type'] == MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER) {
                    ////查询小程序海报
                    foreach ($val['resource'] as &$item) {
                        $item['suncode_mediaid'] = $banners[$item['task_id']]['suncode_mediaid'] ?? '';
                        $item['suncode_url']     = $banners[$item['task_id']]['background_url'] ?? '';
                        $item['app_id']          = $appId;
                        //中程序卡片跳转的详情地址
                        $item['applet_card_page'] = '';
                        if (!empty($item['applet_card_cover_url'])) {
                            $item['applet_card_page'] = $this->getIncludeTraceIdPath($item['applet_path'], $traceParams);
                        }
                    }
                }
            }
        }

        return [
            $result,
            $total
        ];
    }

    /**
     * 获取带有trace——id的路径
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/26 17:24
     *
     * @param $appletPath
     *
     * @return string
     */
    public function getIncludeTraceIdPath($appletPath, $traceParams)
    {
        $path = $appletPath;
        if (strpos($appletPath, '?') !== false) {
            $arr         = explode('?', $appletPath);
            $path        = $arr[0];
            $paramsStr   = $arr[1];
            $queryParams = [];
            parse_str($paramsStr, $queryParams);
            if (isset($queryParams['scene']) && !empty($queryParams['scene'])) {
                $decryptData = decryptShortLink($queryParams['scene']);
                $decryptData = $decryptData['args'] ?? [];
                $queryParams = [];
                foreach ($decryptData as $key => $value) {
                    $queryParams[$value['name']] = $value['value'];
                }
            }
            $traceParams = array_merge($traceParams, $queryParams);
        }
        $traceId = encryptShortLink($traceParams);
        $path    = $path . '.html?short_param_id=' . $traceId;

        return $path;
    }

    /**
     * 获取单个导购的计划可见范围
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/13 16:39
     *
     * @param $enterpriseId
     * @param $workerId
     * @param $qyDomain
     * @param $startTime
     * @param $endTime
     *
     * @throws \Exception
     */
    public function getRangePlanList($enterpriseId, $workerId)
    {
        $rangeList = MarketingPlanVisibleRange::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('worker_id', $workerId)
            ->get()
            ->toArray();

        $planIdArr = array_column($rangeList, 'plan_id');
        $planIdArr = array_unique($planIdArr);

        return $planIdArr;
    }


    /**
     * 格式化每日的数据
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 17:44
     *
     * @param      $data
     * @param bool $returnList
     *
     * @return array
     */
    public function formatDayData(&$taskList, $startTime, $endTime, $currentTime, bool $returnList = true)
    {
        //将时间格式化为日
        $currentTime = strtotime(date('Y-m-d', empty($currentTime) ? time() : $currentTime));
        $startTime   = strtotime(date('Y-m-d', $startTime));
        $endTime     = strtotime(date('Y-m-d', $endTime));
        if (empty($startTime)) {
            $startTime = strtotime(date('Y-m-d'));
        }
        if (empty($endTime)) {
            $endTime = strtotime(date('Y-m-d'));
        }

        $result = [];

        $aDaySeconds = 86400;
        foreach ($taskList as $k => &$item) {
            $loopStartTime = $item['start_time'];
            $loopEndTime   = $item['end_time'];

            if ($loopStartTime < $startTime) {
                $loopStartTime = $startTime;
            }
            if ($loopEndTime > $endTime) {
                $loopEndTime = $endTime;
            }
            while ($loopStartTime <= $loopEndTime) {
                $dateIndex    = date('Y-m-d', $loopStartTime);
                $taskCount    = isset($result[$dateIndex]['task_count']) ? ($result[$dateIndex]['task_count'] + 1) : 1;
                $itemTaskList = [];
                if ($returnList) {
                    if (isset($result[$dateIndex]['task_list'])) {
                        $itemTaskList = $result[$dateIndex]['task_list'];
                    }
                    $itemTaskList[] = $item;
                }
                $itemArr = [
                    'year'       => date('Y', $loopStartTime),
                    'month'      => date('n', $loopStartTime),
                    'day'        => date('j', $loopStartTime),
                    'week_day'   => date('w', $loopStartTime),
                    'is_current' => ($loopStartTime == $currentTime) ? true : false,
                    'display'    => true,
                    'task_count' => $taskCount,
                    'task_list'  => $itemTaskList
                ];

                $result[$dateIndex] = $itemArr;

                $loopStartTime += $aDaySeconds;
            }
        }
        $beginTime = $startTime;
        //组装成时间范围的数据,将初始日期设置为周日开始
        while ($beginTime <= $endTime) {
            $date = date('Y-m-d', $beginTime);
            if (!isset($result[$date])) {
                $itemArr       = [
                    'year'       => date('Y', $beginTime),
                    'month'      => date('n', $beginTime),
                    'day'        => date('j', $beginTime),
                    'week_day'   => date('w', $beginTime),
                    'is_current' => ($beginTime == $currentTime) ? true : false,
                    'display'    => ($beginTime >= $startTime) ? true : false,
                    'task_count' => 0,
                    'task_list'  => []
                ];
                $result[$date] = $itemArr;
            }
            //格式化参数
            $result[$date]['week_day'] = intval($result[$date]['week_day']);

            $beginTime += $aDaySeconds;
        }
        ksort($result);

        return $result;
    }

    /**
     * 组装成日历的完整数据
     *
     * @param $dayDataList [Y-m-d] => [] | 或者key为数字，确保每天的数据可按key正常从前往后进行排序
     *                     User: Hui Yang
     *                     Email: 2361313833@qq.com
     *                     DateTime: 2021/1/13 10:28
     */
    public function fullFormatCalendar(&$dayDataList)
    {
        $result       = [];
        $weekDayCount = 0;
        $weekCount    = 0;
        $weekDataArr  = [];
        foreach ($dayDataList as $k => &$item) {
            $year         = $item['year'];
            $month        = $item['month'];
            $yearMonthKey = $year . '-' . $month;
            //如果当前结果没有该年月，创建
            if (!isset($result[$yearMonthKey])) {
                $result[$yearMonthKey] = [
                    'year'  => $year,
                    'month' => $month,
                    'days'  => []
                ];
                //新的月，周进行初始化
                $weekCount    = 0;
                $weekDayCount = 0;
            }
            //周
            $weekDataArr[] = $item;
            ++$weekDayCount;
            //判断第一个元素是否星期日，如果不是则进行处理
            if ($weekCount == 0 && $weekDayCount == 1) {
                $firstItem = array_first($weekDataArr);

                $firstWeekDay = $firstItem['week_day'];
                if ($firstWeekDay != 0) {
                    for ($i = $firstWeekDay; $i > 0; $i--) {
                        $firstItem['day']        = '';
                        $firstItem['week_day']   = $i - 1;
                        $firstItem['display']    = false;
                        $firstItem['is_current'] = false;
                        $firstItem['task_count'] = 0;
                        $firstItem['task_list']  = [];
                        array_unshift($weekDataArr, $firstItem);
                        $weekDayCount++;
                    }
                }
            }

            //如果星期计数等于7，则初始为0，周计数加1
            if ($weekDayCount >= 7) {
                $weekDayCount = 0;
                ++$weekCount;
            }
            //如果数组的下一个月开始为1，如果已经到了数组的最后一个，则赋值周数据
            $nextDateKey  = date('Y-m-d', strtotime("{$k} +1 day"));
            $lastArrayKey = array_key_last($dayDataList);
            if ($k == $lastArrayKey ||
                (isset($dayDataList[$nextDateKey]['day']) && $dayDataList[$nextDateKey]['day'] == '01')
            ) {
                $result[$yearMonthKey]['days'] = array_values($weekDataArr);
                $weekDataArr                   = [];
            }
        }
        $result = array_values($result);

        return $result;
    }

    /**
     * 未完成人名单
     *
     * @param $params
     *
     * @return mixed
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/15 5:45
     */
    public function notFinished($params)
    {
        $planInfo   = MarketingPlanTask::query()
            ->with(['range' => function ($query) {
            }])
            ->where('id', $params['task_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $planInfo   = !empty($planInfo) ? $planInfo->toArray() : [];
        $objectType = $planInfo['object_type'] ?? 1;
        $workerIds  = array_column($planInfo['range'], 'worker_id');
        $this->getGuideInfoByIds($workerIds);
        //导购任务完成数据
        $this->getCompleteInfoByTaskIds([$params['task_id']]);
        $effectiveGuideIDs = $this->calculationTaskCompletion($params['task_id'], $workerIds);

        //下载获得数据
        //未完成:已完成和全部求差集
        $notFinishedIDS  = array_diff($workerIds, $effectiveGuideIDs);
        $notFinishedInfo = [];
        foreach ($this->guideInfo as $k => $v) {
            if (in_array($v['worker_id'], $notFinishedIDS)) {
                $notFinishedInfo[] = $v;
            }
        }
        //门店范围 获得门店名称
        if ($notFinishedInfo && $objectType == MarketingPlan::OBJECT_TYPE_STORE) {
            $storeIds  = array_values(array_unique(array_column($notFinishedInfo, 'store_id')));
            $storeList = (new Store())->listByIds(['store_ids' => $storeIds, 'is_store' => 1]);
            $storeList = array_column($storeList, 'short_name', 'id');
            foreach ($notFinishedInfo as &$value) {
                $value['store_name'] = $storeList[$value['store_id']] ?? '';
            }
        }
        //部门
        if ($notFinishedInfo && $objectType == MarketingPlan::OBJECT_TYPE_DEPARTMENT) {

            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']]);
            $memUids        = array_column($notFinishedInfo, 'mem_uid');
            // 参数列表
            $params = [
                'memUids'  => $memUids,
                'pageNum'  => Constant::DEFAULT_PAGE,
                'pageSize' => count($memUids)
            ];
            // 调用uc sdk获取导购信息
            $members = (new Member())->list($enterpriseInfo['data']['domain'] ?? '', $enterpriseInfo['data']['identifier'] ?? '', $params);
            $members = array_column($members['list'], null, 'memUid');
            foreach ($notFinishedInfo as &$val) {
                $val['store_name'] = $members[$val['mem_uid']]['dpName'][0]['dpName'] ?? '';
            }

        }
        $data['not_finished'] = $notFinishedInfo;
        $data['object_type']  = $objectType;
        return $data;
    }

    /**
     * 获得门店或部门的名称
     *
     * @param $Ids
     * @param $type
     * @param $params
     *
     * @return array
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/19 19:16
     */
    public function getObjectNameByIds($workerIds)
    {
        $rangeInfo = [];
        $guideInfo = array_column($this->guideInfo, null, 'worker_id');
        foreach ($workerIds as $workerId) {
            if (isset($guideInfo[$workerId])) {
                $rangeInfo[] = [
                    'worker_id' => $guideInfo[$workerId]['worker_id'],
                    'mem_uid'   => $guideInfo[$workerId]['mem_uid'],
                    'name'      => $guideInfo[$workerId]['name'],
                    'worker_no' => $guideInfo[$workerId]['worker_no']
                ];
            }
        }
        return $rangeInfo;
    }

    /**
     * 本地
     *
     * @param string $hz
     *
     * @return string
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/20 16:15
     */
    public function getLocalImgPath($hz = 'png')
    {
        return storage_path('logs/tempimg/') . date("YmdHis") . rand(1000, 9999) . "." . $hz;
    }

    /**
     * @param $params
     *
     * @return bool
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/22 10:23
     */
    public function getAllDepartmentList($params)
    {
        $qyDomain      = $params['qy_domain'] ?? '';
        $appIdentifier = $params['app_identifier'] ?? '';
        //$params['enterprise']  = $params['enterprise_hash'];
        $params['pageNum']  = 1;
        $params['pageSize'] = 9999;
        //查询子目录
        $result = (new Department())->list($qyDomain, $appIdentifier, $params);
        foreach ($result['list'] as $k => $v) {
            $rangeInfo[] = [
                'object_id'   => $v['dpId'],
                'object_name' => $v['dpName']
            ];
        }
        $this->allDepartmentList = $rangeInfo;

        return true;
    }

    /**
     *
     * @param $params
     *
     * @return bool
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/2/2 9:44
     */
    public function getStoreListByIds($params)
    {
        $this->storeList = (new Store())->listByIds(['store_ids' => $params['storeIds'], 'is_store' => 1]);
        return true;
    }

    /**
     * 生成后台预览海报图
     *
     * @param $resource
     * @param $enterpriseId
     *
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/25 11:20
     */
    public function createPreviewPoster($resource)
    {
        //太阳码
        $sunCode = $resource['sun_code'];
        if (empty($resource['sun_code'])) {
            [$sunCode, $qrCodeParams] = $this->getSunCode($resource['applet_path'], [], $resource['enterprise_id'], $resource['task_id']);
        }

        if (!empty($sunCode)) {
            //海报
            $posterUrl              = $this->getPosterImg($sunCode, $resource['resource_url']);
            $resource['poster_url'] = $posterUrl;
            //存储
            MarketingPlanTaskResource::query()
                ->where('id', $resource['resource_id'])
                ->update([
                    'sun_code'   => $sunCode,
                    'poster_url' => $posterUrl,
                    'updated_at' => time()
                ]);
        }
        if (empty($sunCode)) {
            $posterUrl = $resource['resource_url'];
        }

        return [$posterUrl, $sunCode];
    }

    /**
     * 获取任务详情
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/25 20:54
     *
     * @param $params
     *
     * @return array
     */
    public function getTaskDetail($params)
    {
        //将结束时间定义为当天的最后一秒
        $params['end_time'] = date('Y-m-d 23:59:59', $params['end_time']);
        $params['end_time'] = strtotime($params['end_time']);

        $result = MarketingPlanTask::query()
            ->with([
                'resource' => function ($query) use ($params) {
                    $query->select([
                        'id',
                        'task_id',
                        'type',
                        'mediaid',
                        'resource_url',
                        'applet_path',
                        'applet_card_title',
                        'applet_card_cover_url',
                        'source_id'
                    ])->orderBy('id', 'asc');
                }])
            ->where('enterprise_id', $params['enterprise_id'])
            //->whereIn('marketing_plan_task.plan_id', $params['plan_ids'])
            ->where('id', $params['task_id'])
            ->orderBy('created_at', 'DESC')
            ->orderBy('id', 'DESC')
            ->select([
                'id',
                'plan_id',
                'task_title',
                'start_time',
                'end_time',
                'content',
                'type as task_type'
            ])
            ->first();
        if (!isset($result->id)) {
            return [];
        }
        $result = $result->toArray();

        $banners = MemberBannerRelation::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('task_id', $result['id'])
            ->where('worker_id', $params['worker_id'])
            ->get([
                'suncode_mediaid',
                'background_url',
                'task_id'
            ])->toArray();
        $banners = array_column($banners, null, 'task_id');

        // 获取任务状态
        $taskFinished = MarketingPlanTaskProgress::query()
            ->where('task_id', $params['task_id'])
            ->where('worker_id', $params['worker_id'])
            ->get()
            ->toArray();

        $result['is_finished'] = $taskFinished ? MarketingPlanTaskProgress::TASK_ALREADY_FINISHED : MarketingPlanTaskProgress::TASK_NOT_FINISHED;

        $result['tags'] = [];
        if ($result['task_type'] == MarketingPlanTask::TYPE_GROUP_TASK) {
            // 获取任务关联的客户标签
            $tagIds = MarketingTaskTagRelation::query()
                ->where('task_id', $params['task_id'])
                ->where('deleted_at', 0)
                ->pluck('tag_id')
                ->toArray();

            if ($tagIds) {
                $result['tags'] = WorkExternalTag::query()
                    ->whereIn('id', $tagIds)
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->where('deleted_at', 0)
                    ->get(['id', 'name'])
                    ->toArray();
            }
        }

        //小程序素材信息
        $miniSourceIds     = collect($result['resource'])
            ->where('type', MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER)
            ->pluck('source_id')
            ->all();
        $miniTitles        = SourceMiniProgram::query()->whereIn('id', $miniSourceIds)->get();
        $miniIdToPath      = $miniTitles->pluck('path', 'id')->all();
        $miniIdToTitle     = $miniTitles->pluck('title', 'id')->all();
        $miniIdToProgramId = $miniTitles->pluck('mini_program_id', 'id')->all();
        $miniList          = MiniProgram::query()->whereIn('id', $miniTitles->pluck('mini_program_id'))->get()->toArray();
        $miniList          = array_column($miniList, null, 'id');

        //这个appid是发小程序卡片要用的
//        $shopDetail = (new MiniProgram())->detail([
//            'enterprise_id' => $params['enterprise_id']
//        ]);
//        $appId      = $shopDetail['data']['app_id'] ?? '';
        //todo 项目不完整兼容
        $appId                    = '';
        $traceParams['mem_uid']   = $params['mem_uid'] ?? '';
        $traceParams['share_uid'] = $params['mem_uid'] ?? '';
        $traceParams['source']    = self::SHARE_TYPE;
        $traceParams['worker_id'] = $params['worker_id'] ?? 0;
        $traceParams['store_id']  = $params['store_id'];

        //资源的类型只会有一种，将type加在外面，方便前端使用
        $result['type'] = 0;
        if (!empty($result['resource'])) {
            $firstItem      = array_first($result['resource']);
            $result['type'] = $firstItem['type'];

            if ($result['type'] == MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER) {
                ////查询小程序海报
                foreach ($result['resource'] as &$item) {
                    $item['title']           = $miniIdToTitle[$item['source_id']] ?? '';
                    $item['path']            = $miniIdToPath[$item['source_id']] ?? '';
                    $item['app_name']        = $miniList[$miniIdToProgramId[$item['source_id']]]['app_name'] ?? '';
                    $item['app_id']          = $miniList[$miniIdToProgramId[$item['source_id']]]['app_id'] ?? '';
                    $item['origin_id']       = $miniList[$miniIdToProgramId[$item['source_id']]]['origin_id'] ?? '';
                    $item['logo']            = $miniList[$miniIdToProgramId[$item['source_id']]]['logo'] ?? '';
                    $item['suncode_mediaid'] = $banners[$item['task_id']]['suncode_mediaid'] ?? '';
                    $item['suncode_url']     = $banners[$item['task_id']]['background_url'] ?? '';
                    //$item['app_id']          = $appId;
                    //中程序卡片跳转的详情地址
                    $item['applet_card_page'] = '';
                    if (!empty($item['applet_card_cover_url'])) {
                        $item['applet_card_page'] = $this->getIncludeTraceIdPath($item['applet_path'], $traceParams);
                    }
                }
            }
        }

        return $result;
    }

    /**
     * 检查是否有太阳码，没有则生成
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/22 18:12
     *
     * @param $enterpriseId
     * @param $taskId
     * @param $workerId
     *
     * @return bool|array
     */
    public function checkSuncodeBanner($enterpriseId, $taskId, $workerInfo)
    {
        $resource = MarketingPlanTaskResource::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('task_id', $taskId)
            ->where('type', MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER)
            ->first();
        if (!isset($resource->id)) {
            return [];
        }
        //检查banner
        $banner = MemberBannerRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('task_id', $taskId)
            ->where('worker_id', $workerInfo['id'])
            ->first();
        //如果海报已经 存在
        if (isset($banner->background_url)
            && !empty($banner->background_url)
            && !empty($banner->suncode_mediaid)
        ) {
            return $banner->toArray();
        }

        //导购信息
        $enterpriseId = $resource->enterprise_id;

        [$sunCode, $qrCodeParams] = $this->getSunCode(
            $resource->applet_path,
            $workerInfo,
            $enterpriseId,
            $taskId
        );
        if (empty($sunCode)) {
            return [];
        }

        //查询企业信息
        try {
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        } catch (\Exception $e) {
            $enterpriseInfo = [];
        }
        $qyDomain = $enterpriseInfo['domain'] ?? '';
        //生成海报
        $poster  = $this->getPosterImg($sunCode, $resource->resource_url);
        $mediaId = '';
        if (!empty($qyDomain)) {
            try {
                $mediaId = $this->getMediaId(
                    $poster,
                    $qyDomain,
                    Constant::OPERATE_TYPE_SEND
                );
            } catch (\Exception $e) {
                $mediaId = '';
            }
        }

        if (isset($banner->id)) {
            $updateData = [
                'suncode_url'     => $sunCode,
                'suncode_param'   => $qrCodeParams,
                'background_url'  => $poster ?? '',
                'suncode_mediaid' => $mediaId,
                'updated_at'      => time()
            ];
            $banner->update($updateData);
            $result = $banner->toArray();
        } else {
            $insertData       = [
                'enterprise_id'   => $enterpriseId,
                'task_id'         => $taskId,
                'mem_uid'         => $workerInfo['mem_uid'],
                'worker_id'       => $workerInfo['id'],
                'suncode_mediaid' => $mediaId,
                'suncode_url'     => $sunCode,
                'suncode_param'   => $qrCodeParams,
                'background_url'  => $poster ?? '',
                'created_at'      => time(),
                'updated_at'      => time()
            ];
            $insertData['id'] = MemberBannerRelation::query()->insertGetId($insertData);
            $result           = $insertData;
        }

        return $result;
    }

    /**
     * 检查是否有mediaId，没有则生成
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/22 18:12
     *
     * @param $enterpriseId
     * @param $taskId
     * @param $workerId
     *
     * @return bool|array
     */
    public function checkResourceMediaId($enterpriseId, $resourceId, $workerInfo)
    {
        $resource = MarketingPlanTaskResource::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('id', $resourceId)
            ->where('mediaid', '')
            ->first();
        if (!isset($resource->id)) {
            return [];
        }
        //导购信息
        $enterpriseId = $resource->enterprise_id;
        try {
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        } catch (\Exception $e) {
            return [];
        }
        $qyDomain = $enterpriseInfo['domain'];

        if (!empty($qyDomain)) {
            try {
                $mediaId                      = $this->getMediaId(
                    $resource->resource_url,
                    $qyDomain,
                    Constant::OPERATE_TYPE_SEND
                );
                $resource->mediaid            = $mediaId;
                $resource->mediaid_updated_at = time();
                $resource->save();
            } catch (\Exception $e) {
                return [];
            }
        }

        return $resource->toArray();
    }

    /**
     * 获得导购信息
     *
     * @param $workerIds
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/29 14:45
     */
    public function getGuideInfoByIds($workerIds)
    {
        $info = StoreWorker::query()
            ->whereIn('id', $workerIds)
            ->where('worker_status', StoreWorker::IN_SERVICE)
            ->select(['id as worker_id', 'mem_uid', 'name', 'worker_no', 'store_id', 'mobile'])
            ->get()
            ->toArray();
        //$info = array_column($info,null,'worker_id');
        $this->workerIds = array_column($info, 'worker_id');
        $this->guideInfo = $info;
        return true;
    }

    /**
     * 发送提醒
     * @param $taskId
     * @Date: 2021/5/20
     * @Time: 17:29
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function sendMessage($taskId)
    {
        $task = MarketingPlanTask::query()
            ->where('id', $taskId)
            ->first();

        if (!$task) {
            Log::info("任务不存在");
            return;
        }

        $task = $task->toArray();
        $time = time();

        $jobData = [
            'id' => $taskId
        ];

        foreach (['before_start_time', 'before_end_time', 'fixed_time'] as $value) {
            if (isset($task[$value]) && $task[$value] && $task[$value] > 0) {
                if (($task[$value] - $time) > 0) {
                    dispatch((new SyncNoticeMarketingPlanTask($jobData))->delay($task[$value] - $time));
                } else {
                    dispatch((new SyncNoticeMarketingPlanTask($jobData)));
                }
            }
        }
    }

    /**
     * 获取某个导购的所有任务id
     * @param $enterpriseId
     * @param $workerId
     * @return array
     * @Date: 2021/5/20
     * @Time: 18:53
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function getRangeTaskList($enterpriseId, $workerId)
    {
        $rangeList = MarketingPlanVisibleRange::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('worker_id', $workerId)
            ->get()
            ->toArray();

        $taskIdArr = array_column($rangeList, 'task_id');
        $taskIdArr = array_unique($taskIdArr);

        return $taskIdArr;
    }

    /**
     * 获取任务数量
     * @param $params
     * @return int
     * @Date: 2021/5/25
     * @Time: 15:07
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function getTaskNumber($params)
    {
        return MarketingPlanTask::query()
                ->where('marketing_plan_task.enterprise_id', $params['enterprise_id'])
                ->whereIn('marketing_plan_task.id', $params['task_ids'])
                ->when(!empty($params['task_type']), function ($query) use ($params) {
                    $query->where('marketing_plan_task.type', $params['task_type']);
                })
                ->when(!empty($params['task_title']), function ($query) use ($params) {
                    $query->where('marketing_plan_task.task_title', 'like', '%' . $params['task_title'] . '%');
                })
                ->when(!empty($params['start_time']) && !empty($params['end_time']), function ($query) use ($params) {
                    return $query->where([
                        ['marketing_plan_task.start_time', "<=", $params['end_time']],
                        ['marketing_plan_task.end_time', ">=", $params['start_time']]
                    ]);
                })
                ->count() ?? 0;
    }


}
