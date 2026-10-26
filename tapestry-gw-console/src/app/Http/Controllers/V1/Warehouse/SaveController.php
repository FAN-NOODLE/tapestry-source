<?php
/**
 * 新增编辑仓库
 * Created by PhpStorm.
 * User pw
 * Date: 2020/4/18
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Warehouse;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Warehouse;
use Illuminate\Support\Facades\Log;

class SaveController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        // 发布验证
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['operator_id']   = $this->enterprise['id'] ?? 0;
        $params['operator_name'] = $this->enterprise['contact'] ?? '';

        try {
            (new Warehouse())->save($params);
        } catch (\Exception $exception) {
            Log::info('仓库操作失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }

        return success([]);
    }
}
