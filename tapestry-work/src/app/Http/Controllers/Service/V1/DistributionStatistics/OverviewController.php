<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\DistributionStatistics;


use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionStatisticsService;
use App\Services\GuideManagerService;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;

class OverviewController extends BaseController
{
    // 排行榜按数量排序类型
    const RANKING_SORT_TYPE = 2;


    /**
     **
     * StatisticsConsole
     *
     * @desc 管理后台顾客概览
     * @return array
     * @date 2020-08-22
     * @throws \Exception
     */
    public function statisticsConsole()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');
        $guideIdList  = $this->request->input('worker_ids');
        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList
                ], false)['list'] ?? [];

        } else {
            $guideList = [];
        }

        $result = (new DistributionStatisticsService())->statisticData($enterpriseId, $guideList, $startTime, $endTime);

        return success($result);
    }

    /**
     * statisticsByDate
     *
     * @desc 按照日期获取分销数据概览
     * @return array
     * @date 2020-08-22
     * @throws \Exception
     */
    public function statisticsByDate()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');;
        $guideIdList     = $this->request->input('worker_ids', []);
        $page            = $this->request->input('page', Constant::PAGE);
        $size            = $this->request->input('size', Constant::SIZE);
        $sortColumn      = $this->request->input('sort_type');
        $statisticalSort = $this->request->input('statistical_sort');
        $forPage         = $this->request->input('for_page', true);

        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList
                ], false)['list'] ?? [];
        } else {
            $guideList = [];
        }

        $result     = (new DistributionStatisticsService())->statisticData($enterpriseId, $guideList, $startTime, $endTime);

        $result     = $result['detail_data'];
        $timeColumn = $columnData=[];
        if (!empty($sortColumn)) {
            foreach ($result as $key => $value) {
                $columnData[$key] = $value[$sortColumn];
                $timeColumn[$key] = $value['time'];
            }

            if ($statisticalSort == 'desc') {
                array_multisort($columnData, SORT_DESC, $timeColumn, SORT_DESC, $result);
            } else {
                array_multisort($columnData, SORT_ASC, $timeColumn, SORT_DESC, $result);
            }
        } else {
            foreach ($result as $key => $value) {
                $timeColumn[$key] = $value['time'];
            }
            array_multisort($timeColumn, SORT_DESC, $result);
        }

        $total = count($result);

        if ($forPage) {
            $result = array_chunk($result, $size)[$page - 1] ?? [];
        }

        return success([
            'total' => $total,
            'page'  => $page,
            'size'  => $size,
            'list'  => $result,
        ]);
    }

    /**
     * statisticsByGuide
     *
     * @desc 根据导购获取分销数据
     * @return array
     * @date 2020-08-23
     * @throws \Exception
     */
    public function statisticsByGuide()
    {

        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');;
        $guideIdList     = $this->request->input('worker_ids', []);
        $page            = $this->request->input('page', Constant::PAGE);
        $size            = $this->request->input('size', Constant::SIZE);
        $sortColumn      = $this->request->input('sort_type');
        $statisticalSort = $this->request->input('statistical_sort');
        $qyDomain        = $this->request->input('qy_domain');
        $forPage         = $this->request->input('for_page', true);

        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList
                ], false)['list'] ?? [];
        } else {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId
                ], false)['list'] ?? [];
        }

        $result = (new DistributionStatisticsService())->statisticDataByGuide($enterpriseId, $guideList, $startTime, $endTime);

        $total = count($result);
        if (!empty($sortColumn)) {
            foreach ($result as $key => $value) {
                $columnData[$key] = $value[$sortColumn];
            }
            if ($statisticalSort == 'desc') {
                array_multisort($columnData, SORT_DESC, $result);
            } else {
                array_multisort($columnData, SORT_ASC, $result);
            }
        }

        if ($forPage) {
            $result = array_chunk($result, $size)[$page - 1] ?? [];
        }

        $result = $this->formatGuideInfo($result, $qyDomain);
        return success([
            'total' => $total,
            'page'  => $page,
            'size'  => $size,
            'list'  => $result,
        ]);
    }


    /**
     * formatGuideInfo
     *
     * @desc 格式化导购信息
     *
     * @param $guideList
     * @param $qyDomain
     *
     * @return mixed
     * @date 2020-08-23
     * @throws \Exception
     */
    private function formatGuideInfo($guideList, $qyDomain)
    {
        $storeIdList = array_unique(array_column($guideList, 'store_id'));

//        $storeList = (new Store())->listByIds(['store_ids' => $storeIdList]);
//        $storeList = array_column($storeList, null, 'id');
        $storeList = [];
        $memUids   = array_column($guideList, 'mem_uid');

        $guideDepartmentInfo = GuideManagerService::getGuideDepartment($memUids, $qyDomain);
        $guideDepartmentInfo = array_column($guideDepartmentInfo, null, 'memUid');
        array_walk($guideList, function (&$value) use ($guideDepartmentInfo, $storeList) {
            $value['dp_name']    = $guideDepartmentInfo[$value['mem_uid']]['dpName'] ?? '';
            $value['store_name'] = '';
            if (!empty($storeList[$value['store_id']])) {
                $value['store_name'] = $storeList[$value['store_id']]['short_name'] ?: $storeList[$value['store_id']]['name'];
            }
            $value['worker_name'] = $value['name'];
        });
        return $guideList;
    }

    /**
     * 智慧导购统计首页
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public function distributionStatistic()
    {

        $param    = $this->request->input();
        $rules    = [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'qy_domain'     => 'required',
            'admin_id'      => 'required',
            'start_time'    => 'required',
            'end_time'      => 'required',
        ];
        $messages = [
            'enterprise_id.*'     => '企业id参数不合法',
            'store_id.*'          => '门店id参数不合法',
            'mobile.*'            => '用户手机号参数不合法',
            'qy_domain.required'  => 'qy_domain参数缺失',
            'admin_id.required'   => 'admin_id参数缺失',
            'start_time.required' => 'start_time参数缺失',
            'end_time.required'   => 'end_time参数缺失',
        ];
        $this->validate($this->request, $rules, $messages);

        $qyDomain     = $param['qy_domain'];
        $adminId      = $param['admin_id'];
        $mobile       = $param['mobile'];
        $enterpriseId = $param['enterprise_id'];
        $storeId      = $param['store_id'];

        switchAppHost(['enterprise_id' => $param['enterprise_id']]);

        // 获取导购列表
        $guideList = GuideManagerService::getGuideListForStatistic($param, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId);

        $result = (new DistributionStatisticsService())->statisticData($enterpriseId, $guideList['guide_list'], $param['start_time'], $param['end_time']);

        return success($result);
    }

    /**
     * 导购排行榜
     */
    public function statisticRanking()
    {

        $input    = $this->request->input();
        $rules    = [
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'type'          => 'required|integer|in:1,2',
            'qy_domain'     => 'required',
            'mobile'        => 'required',
            'enterprise_id' => 'required|integer',
            'store_id'      => 'required|integer',
            'admin_id'      => 'required',
        ];
        $messages = [
            'type.required'          => 'type不能为空',
            'qy_domain.required'     => '企业domain不能为空',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID必须为正整数',
            'store_id.required'      => '门店id不能为空',
            'admin_id.required'      => 'admin_id不能为空',
        ];
        $this->validate($this->request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);

        // 获取对应导购列表
        $guideList = GuideManagerService::getGuideListForRank($input['mobile'], $input['qy_domain'], $input['enterprise_id'], $input['admin_id'], $input['store_id']);
        if (empty($guideList['guide_list'])) {
            return success();
        }

        $result = (new DistributionStatisticsService())->statisticDataByGuide($input['enterprise_id'], $guideList['guide_list'], $input['start_time'], $input['end_time']);

        $sortColumn = 'settlement_money';
        // 按数量排行
        if ($input['type'] == self::RANKING_SORT_TYPE) {
            $sortColumn = 'settlement_num';
        }

        foreach ($result as $key => $value) {
            $columnData[$key] = $value[$sortColumn];
        }

        array_multisort($columnData, SORT_DESC, $result);
        $result = array_slice($result, 0, 50);

        $result = $this->formatGuideInfo($result, $input['qy_domain']);

        return success($result);
    }


}
