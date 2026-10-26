<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;


use App\Common\DateUtils;
use App\Models\StoreDepartment;
use App\Models\WorkerLevelArea;
use App\Services\DataStatistic\CustomerDataCompassService;
use App\Services\EnterpriseService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\EnterpriseConfig;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\Rule;
use App\Common\Constant;

class RankForConsoleController extends BaseController
{
    const LEVEL_KEY = 'area_level_key';


    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'worker_ids' => 'array',
            'qy_domain' => 'required|string',
            'admin_id' => 'required',
            'start_time' => 'required|integer',
            'end_time' => 'required|integer',
            'rank_column' => [
                'required',
                'string',
                Rule::in(array_keys(CustomerDataCompassService::DEFAULT_STATISTIC))
            ],
            'rank_length' => 'required|integer'
        ], [
            'enterprise_id.*' => '企业id参数不合法',
            'qy_domain.required' => 'qy_domain参数不合法',
            'admin_id.required' => 'admin_id参数缺失',
            'start_time.*' => 'start_time参数不合法',
            'end_time.*' => 'end_time参数不合法',
            'rank_column.*' => 'rank_column参数不合法',
            'rank_length.*' => 'rank_length参数不合法',
        ]);
        $enterpriseId = $this->request->input('enterprise_id');
        $qyDomain = $this->request->input('qy_domain');
//        $adminId = $this->request->input('admin_id', 0);
        $rankColumn = $this->request->input('rank_column');
        $rankLength = $this->request->input('rank_length');
        $startTime = $this->request->input('start_time');
        $endTime = $this->request->input('end_time');
        $startStatTime = DateUtils::timeToStatTime($startTime);
        $endStatTime = $endTime;

//        switchAppHost(['enterprise_id' => $enterpriseId]);

        $result = StoreWorkerService::list(['enterprise_id' => $enterpriseId, 'worker_status' => Constant::IN_SERVICE], false);
        $workerList = $result['list'] ?? [];
        if (!empty($workerList)) {
            $workerList = array_values(array_column($workerList, null, 'mem_uid'));
            $rankStatisticList = $this->customerDataCompassService->rankForConsole($enterpriseId, $startStatTime, $endStatTime,
                 $rankLength, $workerList);
            $this->customerDataCompassService->decorateWorkerInfo($rankStatisticList, $workerList, $qyDomain);
        } else {
            Log::info('customer statistic rank, workerList is empty');
            $rankStatisticList = [];
        }
        //获取企业导购区域配置
        $configList = EnterpriseService::getEnterpriseConfig($enterpriseId, self::LEVEL_KEY);
        $enterpriseConfig = $configList['data'] ?? [];

        $areaLevel = 0;
        if (!empty($enterpriseConfig)) {
            $areaLevel = current(array_column($enterpriseConfig, 'value'));
            if ($areaLevel != 0) {
                foreach ($rankStatisticList as &$val) {
                    $pathArr = explode('/', $val['dp_name_path']);
                    $val['area_level_list']  = array_slice($pathArr, 0-$areaLevel);
                }
            }
        }
        return success([
            'area_level_config' => $areaLevel,
            'list' => $rankStatisticList

        ]);
    }
}