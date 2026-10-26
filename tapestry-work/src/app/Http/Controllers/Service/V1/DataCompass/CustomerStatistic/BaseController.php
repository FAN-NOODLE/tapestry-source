<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;


use App\Common\Constant;
use App\Common\DateUtils;
use App\Exceptions\Error;
use App\Models\StoreWorker;
use App\Services\DataStatistic\CustomerDataCompassService;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\Rule;

class BaseController extends \App\Http\Controllers\Service\V1\BaseController
{
    /**
     * @var CustomerDataCompassService
     */
    protected $customerDataCompassService;


    public function __construct(CustomerDataCompassService $customerDataCompassService)
    {
        parent::__construct();
        $this->customerDataCompassService = $customerDataCompassService;
    }

    protected function getAdminListRequestParams()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'qy_domain' => 'required|string',
            'start_time' => 'required|integer',
            'end_time' => 'required|integer',
            'worker_ids' => 'array',
            'sort_column' => [
                'string',
                Rule::in(array_keys(CustomerDataCompassService::DEFAULT_STATISTIC))
            ],
            'sort_type' => 'string',
            'page' => 'integer',
            'size' => 'integer',
            'for_page' => 'boolean',
        ];
        $this->validate($this->request, $rules);

        $params['enterprise_id'] = $this->request->input('enterprise_id');
        $params['qy_domain'] = $this->request->input('qy_domain');
        $params['start_time'] = $this->request->input('start_time');
        $params['end_time'] = $this->request->input('end_time');
        $params['worker_ids'] = $this->request->input('worker_ids', []);
        $params['sort_column'] = $this->request->input('sort_column', '');
        $params['sort_type'] = $this->getSortType($this->request->input('sort_type', 'desc'));
        $params['page'] = $this->request->input('page', Constant::PAGE);
        $params['size'] = $this->request->input('size', Constant::SIZE);
        $params['for_page'] = $this->request->input('for_page', true);
        $params['start_stat_time'] = DateUtils::timeToStatTime($params['start_time']);
        $params['end_stat_time'] = DateUtils::timeToStatTime($params['end_time']);
        return $params;
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
                'code' => $e->getCode(),
                'statusCode' => $e instanceof Error ? $e->getStatusCode() : 0,
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            $workerList = [];
        }
        return $workerList;
    }

    protected function getSortType(string $sortType): int
    {
        switch ($sortType) {
            case 'asc':
                return SORT_ASC;
            case 'desc':
            default:
                return SORT_DESC;
        }
    }
}