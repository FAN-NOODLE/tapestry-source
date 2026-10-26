<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/2
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\MaterialRequest;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class UpdateStockController extends BaseController
{
    /**
     * 更新库存--库存占用
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params = $this->request->only([
            'material_id',
            'num', // 库存操作数量
            'operate_type', // 1: 增加 2：减少 3：覆盖
            'source_id', // 来源id // 门店 商城
            'operate_id',
            'is_channel',
            'batch_no',
            'desc'
        ]);

        MaterialRequest::updateStock($params);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['type'] = 1;
        $params['operate_name'] = $this->enterpriseAdmin['username'];
        $params['operate_id'] = $this->enterpriseAdminId;

        // 库存更新
        try {
            (new Material())->updateStock($params);
        } catch (Error $error) {
            Log::info('更新物料库存失败', $error->data);
            notice(451516);
        } catch (\Exception $exception) {
            Log::info('更新物料库存失败'.$exception->getMessage(), $params);
            notice(451516);
        }

        $this->writeOperatorLog("管理员:".$params['operate_name'].", ID：".$params['operate_id'].",更新物料库存");

        return success();
    }

}