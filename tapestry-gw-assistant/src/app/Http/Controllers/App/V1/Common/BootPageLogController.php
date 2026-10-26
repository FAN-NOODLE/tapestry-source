<?php

namespace App\Http\Controllers\App\V1\Common;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\StoreArea;
use CrsSdk\Work\Common;
use Illuminate\Support\Facades\Auth;

class BootPageLogController extends Controller
{
    /**记录引导页浏览记录
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 10:24
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $param['version_code'] = $this->request->post('version_code');

        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'] ?? 0;
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);
        $param['enterprise_id'] = $enterpriseId;
        $param['worker_id'] = $workerDetail['id'] ?? 0;

        $result = (new Common())->bootLogCreate($param);

        return success($result['data'] ?? $result);
    }

    /**
     * 检查引导页是否需要展示
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 10:25
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function check()
    {
        $param = $this->request->input();

        $param['version_code'] = $this->request->post('version_code');

        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $param['enterprise_id'] = $enterpriseId;
        $param['worker_id'] = $workerDetail['id'] ?? 0;

        $result = (new Common())->bootLogCheck($param);

        return success($result['data'] ?? $result);
    }
}
