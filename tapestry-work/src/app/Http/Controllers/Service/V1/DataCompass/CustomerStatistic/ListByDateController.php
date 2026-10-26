<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;


class ListByDateController extends BaseController
{

    public function index()
    {
        $params = $this->getAdminListRequestParams();

        $workerList = [];
        if (!empty($params['worker_ids'])) {
            $workerList = $this->listServicingWorker($params['enterprise_id'], $params['worker_ids']);
            if (empty($workerList)) {
                return [
                    'total' => 0,
                    'page' => $params['page'],
                    'size' => $params['size'],
                    'list' => []
                ];
            }
        }

        $result = $this->customerDataCompassService->listByDate(
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

        return success([
            'total' => $result['total'],
            'page' => $params['page'],
            'size' => $params['size'],
            'list' => $result['list'],
        ]);
    }
}