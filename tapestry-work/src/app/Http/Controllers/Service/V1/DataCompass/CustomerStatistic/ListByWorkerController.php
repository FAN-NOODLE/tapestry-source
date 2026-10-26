<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;



class ListByWorkerController extends BaseController
{

    public function index()
    {
        $params = $this->getAdminListRequestParams();

        $workerList = $this->listServicingWorker($params['enterprise_id'], $params['worker_ids']);
        if (empty($workerList)) {
            return success([
                'total' => 0,
                'page' => $params['page'],
                'size' => $params['size'],
                'list' => [],
            ]);
        }

        $result = $this->customerDataCompassService->listByWorker(
            $params['enterprise_id'],
            $params['start_stat_time'],
            $params['end_stat_time'],
            $workerList,
            $params['sort_column'],
            $params['sort_type'],
            $params['page'],
            $params['size'],
            $params['for_page']
        );

        $this->customerDataCompassService->decorateWorkerInfo($result['list'], $workerList, $params['qy_domain']);

        return success([
            'total' => $result['total'],
            'page' => $params['page'],
            'size' => $params['size'],
            'list' => $result['list'],
        ]);
    }
}