<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/8/17
 * Time: 11:14
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerService;
use App\Services\GuideManagerService;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;

class CustomerStatisticController extends BaseController
{
    /**
     * 小程序顾客数据概览
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $param        = $this->getParam();
        $qyDomain     = $param['qy_domain'];
        $adminId      = $param['admin_id'];
        $mobile       = $param['mobile'];
        $enterpriseId = $param['enterprise_id'];
        $storeId      = $param['store_id'] ?? 0;

        switchAppHost(['enterprise_id' => $enterpriseId]);

        $guideList = GuideManagerService::getGuideListForStatistic($param, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId);
        $result    = CustomerService::statisticCustomerData($enterpriseId, $guideList['guide_list'], $param['start_time'], $param['end_time']);
        return success($result);
    }

    private function getParam()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'qy_domain'     => 'required',
            'admin_id'      => 'required',
            'start_time'    => 'required',
            'end_time'      => 'required',
        ], [
            'enterprise_id.*'     => '企业id参数不合法',
            'store_id.*'          => '门店id参数不合法',
            'mobile.*'            => '用户手机号参数不合法',
            'qy_domain.required'  => 'qy_domain参数缺失',
            'admin_id.required'   => 'admin_id参数缺失',
            'start_time.required' => 'start_time参数缺失',
            'end_time.required'   => 'end_time参数缺失',
        ]);
        return $this->request->only(['enterprise_id', 'store_id', 'mobile', 'qy_domain', 'admin_id', 'start_time', 'end_time',
                                     'worker_ids', 'store_ids', 'department_ids']);
    }

    public function rank()
    {
        $param        = $this->getRankParam();
        $qyDomain     = $param['qy_domain'];
        $adminId      = $param['admin_id'];
        $mobile       = $param['mobile'];
        $enterpriseId = $param['enterprise_id'];
        $storeId      = $param['store_id'] ?? 0;

        switchAppHost(['enterprise_id' => $enterpriseId]);
        $guideList = GuideManagerService::getGuideListForRank($mobile, $qyDomain, $enterpriseId, $adminId, $storeId);
        $result    = GuideManagerService::guideRankOfCustomerData($enterpriseId, $qyDomain, $guideList['guide_list'], $param['type'], $param['start_time'], $param['end_time']);
        return success(['list' => $result]);
    }

    private function getRankParam()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'qy_domain'     => 'required',
            'admin_id'      => 'required',
            'start_time'    => 'required',
            'end_time'      => 'required',
            'type'          => 'required',
        ], [
            'enterprise_id.*'     => '企业id参数不合法',
            'store_id.*'          => '门店id参数不合法',
            'mobile.*'            => '用户手机号参数不合法',
            'qy_domain.required'  => 'qy_domain参数缺失',
            'admin_id.required'   => 'admin_id参数缺失',
            'start_time.required' => 'start_time参数缺失',
            'end_time.required'   => 'end_time参数缺失',
            'type.required'       => 'type参数缺失',
        ]);
        return $this->request->only(['enterprise_id', 'store_id', 'mobile', 'qy_domain', 'admin_id', 'start_time', 'end_time', 'type']);
    }

    /**
     * guideList
     *
     * @desc shop端调用，获取范围内导购列表
     * @return array
     * @date 2020-08-22
     * @throws \Exception
     */
    public function guideList()
    {
        $param        = $this->getGuideListParam();
        $qyDomain     = $param['qy_domain'];
        $adminId      = $param['admin_id'];
        $mobile       = $param['mobile'];
        $enterpriseId = $param['enterprise_id'];
        $storeId      = $param['store_id'] ?? 0;
        $guideList    = [];

        switchAppHost(['enterprise_id' => $enterpriseId]);
        switch ($param['data_type']) {
            case 1:
                $guideList = GuideManagerService::getGuideListForStatistic($param, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId);
                break;
            case 2:
                $guideList = GuideManagerService::getGuideListForRank($mobile, $qyDomain, $enterpriseId, $adminId, $storeId);
                break;
        }
        return success($guideList);
    }

    private function getGuideListParam()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'qy_domain'     => 'required',
            'admin_id'      => 'required',
            'data_type'     => 'required',
        ], [
            'enterprise_id.*'    => '企业id参数不合法',
            'store_id.*'         => '门店id参数不合法',
            'mobile.*'           => '用户手机号参数不合法',
            'qy_domain.required' => 'qy_domain参数缺失',
            'admin_id.required'  => 'admin_id参数缺失',
            'type.required'      => 'type参数缺失',
        ]);
        return $this->request->only(['enterprise_id', 'store_id', 'mobile', 'qy_domain', 'admin_id', 'data_type', 'start_time', 'end_time',
                                     'worker_ids', 'store_ids', 'department_ids']);
    }

    /**
     * customerStatistics4Console
     *
     * @desc 管理后台顾客概览
     * @return array
     * @date 2020-08-22
     * @throws \Exception
     */
    public function customerStatisticsConsole()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');
        $guideIdList  = $this->request->input('worker_ids');
        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList,
                ], false)['list'] ?? [];
        } else {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                ], false)['list'] ?? [];
        }
        $result = CustomerService::statisticCustomerData($enterpriseId, $guideList, $startTime, $endTime);
        return success($result);
    }

    /**
     * customerStatisticsByDate
     *
     * @desc 按照日期获取顾客数据概览
     * @return array
     * @date 2020-08-22
     * @throws \Exception
     */
    public function customerStatisticsByDate()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');;
        $guideIdList = $this->request->input('worker_ids', []);
        $page        = $this->request->input('page', Constant::PAGE);
        $size        = $this->request->input('size', Constant::SIZE);
        $sortColumn  = $this->request->input('sort_column');
        $sortType    = (int)$this->request->input('sort_type');
        $forPage     = $this->request->input('for_page', true);
        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $forAll    = false;
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList,
                ], false)['list'] ?? [];
        } else{
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                ], false)['list'] ?? [];
        }

        $result     = CustomerService::statisticCustomerData($enterpriseId, $guideList, $startTime, $endTime);
        $result     = $result['detail_data'];
        $timeColumn = [];
        if (!empty($sortColumn)) {
            foreach ($result as $key => $value) {
                $columnData[$key] = $value[$sortColumn];
                $timeColumn[$key] = $value['time'];
            }
            array_multisort($columnData, $sortType, $timeColumn, SORT_DESC, $result);
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
     * customerStatisticsByGuide
     *
     * @desc 根据导购获取顾客数据
     * @return array
     * @date 2020-08-23
     * @throws \Exception
     */
    public function customerStatisticsByGuide()
    {

        $enterpriseId = $this->request->input('enterprise_id');
        $startTime    = $this->request->input('start_time');
        $endTime      = $this->request->input('end_time');;
        $guideIdList = $this->request->input('worker_ids', []);
        $page        = $this->request->input('page', Constant::PAGE);
        $size        = $this->request->input('size', Constant::SIZE);
        $sortColumn  = $this->request->input('sort_column');
        $sortType    = (int)$this->request->input('sort_type');
        $qyDomain    = $this->request->input('qy_domain');
        $forPage     = $this->request->input('for_page', true);

        switchAppHost(['enterprise_id' => $enterpriseId]);

        if (!empty($guideIdList)) {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                    'ids'           => $guideIdList,
                ], false)['list'] ?? [];
        } else {
            $guideList = StoreWorkerService::list([
                    'enterprise_id' => $enterpriseId,
                ], false)['list'] ?? [];
        }

        $result = CustomerService::statisticCustomerDataByGuide($enterpriseId, $guideList,$startTime, $endTime);
        $total = count($result);
        if (!empty($sortColumn)) {
            foreach ($result as $key => $value) {
                $columnData[$key] = $value[$sortColumn];
            }
            array_multisort($columnData, $sortType, $result);
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
        $storeIdList         = array_column($guideList, 'store_id');

        $storeList           = (new Store())->listByIds(['store_ids' => $storeIdList]);
        $storeList           = array_column($storeList, null, 'id');
        $memUids             = array_column($guideList, 'mem_uid');

        $guideDepartmentInfo = GuideManagerService::getGuideDepartment($memUids, $qyDomain);
        $guideDepartmentInfo = array_column($guideDepartmentInfo, null, 'memUid');
        array_walk($guideList, function (&$value) use ($guideDepartmentInfo, $storeList) {
            $value['dp_name']          = $guideDepartmentInfo[$value['mem_uid']]['dpName'] ?? '';;
            $value['store_name']       = $storeList[$value['store_id']]['name'] ?? '';
            $value['worker_name']      = $value['name'];
        });
        return $guideList;

    }
}