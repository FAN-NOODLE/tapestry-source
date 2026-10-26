<?php
/*
 * @Author: zhaojie
 * @Date: 2020-03-18 17:45:55
 * @LastEditTime: 2020-03-21 16:14:33
 * @LastEditors: Please set LastEditors
 * @Description: 根据手机号获取门店Ids
 * @FilePath: \gw-shop\src\app\Http\Controllers\Service\V1\Store\Worker\DetailController.php
 */

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class GetWorkMemberInfoController extends BaseController
{
    /**
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();
        $data   = [];
        if (!empty($params)) {
            $data = StoreWorkerService::selectWorkerDetailList($params);
        }
        return success($data);
    }
}
