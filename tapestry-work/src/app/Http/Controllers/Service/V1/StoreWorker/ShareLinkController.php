<?php


namespace App\Http\Controllers\Service\V1\StoreWorker;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class ShareLinkController extends BaseController
{

    /**
     * 导购分享链路参数
     * @return array
     */
    public function index()
    {
        $enterpriseId = $this->request->post('enterprise_id', 0);
        $memUid = $this->request->post('mem_uid', '');
        $sceneUserId = $this->request->post('user_id', 0);
        $storeId = $this->request->post('store_id', 0);

        $result = StoreWorkerService::getWorkerByShareLink($enterpriseId, $memUid, $sceneUserId, $storeId);
        return success($result);
    }
}
