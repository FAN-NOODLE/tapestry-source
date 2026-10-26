<?php


namespace App\Http\Controllers\Service\V1\DataLake;

use Throwable;

class BaBindQueryController extends DataLakeBaseController
{
    private $bindType;

    public function __construct()
    {
        $this->bindType = request()->input('bind_type', "0");

        parent::__construct();
    }

    /**
     * @throws Throwable
     */
    public function requestApi()
    {
        $result = $this->dataLakeService->baBindQuery($this->type, $this->value, $this->bindType);

        // 处理DL空格问题
        foreach ($result as $index => $item) {
            if (isset($item['bindstatus'])) {
                $result[$index]['bindstatus'] = (int)$item['bindstatus'];
            }
            if (isset($item['bindtype'])) {
                $result[$index]['bindtype'] = (int)$item['bindtype'];
            }
        }

        return $result;
    }
}