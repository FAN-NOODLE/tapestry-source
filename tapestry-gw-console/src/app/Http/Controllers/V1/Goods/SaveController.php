<?php
/**
 * 商品新增
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SaveController extends BaseController
{

    /**
     * 新增商品
     *
     * @return mixed
     * @throws \Exception
     * @api v1/material/add
     */
    public function index()
    {
        // 发布验证
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['source_id']     = $this->enterpriseId;
        $params['operator_id']   = $this->enterprise['id'] ?? 0;
        $params['operator_name'] = $this->enterprise['contact'] ?? '';
        $this->writeOperatorLog('保存商品信息');
        try {
            $res = (new Goods())->save($params);
        } catch (\Exception $exception) {
            Log::info('商品操作失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }

        return success($res);
    }
}
