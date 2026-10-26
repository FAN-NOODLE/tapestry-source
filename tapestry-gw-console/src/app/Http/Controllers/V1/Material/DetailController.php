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
use App\Services\MaterialService;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\MicroGoodsV1\WarehouseStock;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class DetailController extends BaseController
{
    /**
     * 获取物料详情--库存物料详情
     *
     * @return mixed
     * @api app/v1/material/detail
     */
    public function index()
    {
        $params = $this->request->all();

        $this->getParams($params);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['source_id'] = $this->request->input('source_id', $this->enterpriseId);

        // 获取物料详情
        try {

            $res = (new Material())->detail($params);

        } catch (Error $error) {

            Log::info('获取详情失败', $error->data);
            notice(450000, $error->getMessage());

        } catch (\Exception $exception) {

            Log::info('获取详情失败' . $exception->getMessage(), $params);
            notice(450000, $exception->getMessage());
        }

        // 获取库存
        if (!empty($res)) {
            $backArr = (new MaterialService())->getTotalStockArray($params['enterprise_id'], $res['id']);
            $res['channel_count'] = (!empty($backArr) ? $backArr['total_channel'] : 0);
            $stockArr = (new WarehouseStock())->getStock([
                'enterprise_id' => $this->enterpriseId,
                'source_id'     => $params['source_id'],
                'material_id'   => $params['material_id']
            ]);
            $res['stock'] = $stockArr['stock'];
            $res['total_stock'] = $stockArr['total_stock'];
        }

        return success($res ?? []);
    }


    protected function getParams($input)
    {
        // 参数验证规则
        $rules = [
            'material_id' => 'required|integer|min:1',
            'source_id'   => 'integer|min:1',
        ];

        // 提示语
        $messages = [
            'material_id.*' => '物料主键必填',
            'source_id.*'   => '门店ID有误',
        ];

        return Validator::make($input, $rules, $messages);
    }
}