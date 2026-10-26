<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;


use App\Common\DateUtils;
use App\Services\DataStatistic\CustomerDataCompassService;
use App\Services\GuideManagerService;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\Rule;

class OverviewForAppController extends BaseController
{
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id' => 'integer',
            'mobile' => 'required',
            'worker_ids' => 'array',
            'qy_domain' => 'required|string',
            'admin_id' => 'required',
            'start_time' => 'required|integer',
            'end_time' => 'required|integer',
        ], [
            'enterprise_id.*' => '企业id参数不合法',
            'store_id.*' => '门店id参数不合法',
            'mobile.*' => '用户手机号参数不合法',
            'qy_domain.required' => 'qy_domain参数不合法',
            'admin_id.required' => 'admin_id参数缺失',
            'start_time.*' => 'start_time参数不合法',
            'end_time.*' => 'end_time参数不合法',
        ]);
        $enterpriseId = $this->request->input('enterprise_id');
        $qyDomain = $this->request->input('qy_domain');
        $mobile = $this->request->input('mobile');
        $storeId = $this->request->input('store_id', 0);
        $adminId = $this->request->input('admin_id', 0);
        $startTime = $this->request->input('start_time');
        $endTime = $this->request->input('end_time');
        $startStatTime = DateUtils::timeToStatTime($startTime);
        $endStatTime = DateUtils::timeToStatTime($endTime);
        $workerIds = $this->request->input('worker_ids');
        $storeIds = $this->request->input('store_ids');
        $departmentIds = $this->request->input('department_ids');

        $workerQueryParams = [
            'worker_ids' => $workerIds,
            'store_ids' => $storeIds,
            'department_ids' => $departmentIds,
        ];

        $workerList = GuideManagerService::getGuideListForStatistic($workerQueryParams, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId);

        if (!empty($workerList['guide_list'])) {
            $workerList = $workerList['guide_list'];
            $overview = $this->customerDataCompassService->overview(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $workerList
            );

            $listByDate = $this->customerDataCompassService->listByDate(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $workerList,
                '',
                0,
                0,
                0,
                false,
                false
            );
        } else {
            Log::info('customer statistic rank, workerList is empty');
            $overview = CustomerDataCompassService::DEFAULT_STATISTIC;
            $listByDate = [];
        }
        return success([
            'overview' => $overview,
            'detail_list' => $listByDate['list']
        ]);
    }
}