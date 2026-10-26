<?php


namespace App\Http\Controllers\Service\Brand\CommunityStatistic;

use App\Exceptions\Error;
use App\Services\DataStatistic\CommunityDataCompassService;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Log;

class BaseController extends \App\Http\Controllers\Service\V1\BaseController
{
    /**
     * @var CommunityDataCompassService
     */
    protected $communityDataCompassService;


    public function __construct(CommunityDataCompassService $communityDataCompassService)
    {
        parent::__construct();
        $this->communityDataCompassService = $communityDataCompassService;
    }

    /**
     * 查询导购列表数据
     * @param int $enterpriseId
     * @param array $workerIdList
     * @param array $columns
     * @return array
     * @throws \Exception
     */
    public function listServicingWorker(int $enterpriseId, array $workerIdList = [], array $columns = ['*'])
    {
        // 查询导购列表
        $workerQueryParams = [
            'enterprise_id' => $enterpriseId,
        ];
        if (!empty($workerIdList)) {
            $workerQueryParams['ids'] = $workerIdList;
        }
        try {
            $workerList = StoreWorkerService::list($workerQueryParams, false, $columns)['list'];
        } catch (\Exception $e) {
            if (env('APP_ENV') != 'prod') {
                throw $e;
            }
            Log::info("导购列表查询异常", [
                'code'       => $e->getCode(),
                'statusCode' => $e instanceof Error ? $e->getStatusCode() : 0,
                'message'    => $e->getMessage(),
                'exception'  => $e->getTraceAsString()
            ]);
            $workerList = [];
        }
        return $workerList;
    }

}