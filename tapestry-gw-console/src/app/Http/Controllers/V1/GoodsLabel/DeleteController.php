<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-07-26
 * Time: 18:56
 */

namespace App\Http\Controllers\V1\GoodsLabel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GoodsLabel;

class DeleteController extends BaseController
{

    /**
     * 删除商品标签
     * {domain}/app/v1/goods-label/delete
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 标签id
        $id = request()->post('id', 0);

        // 标签id不能为空
        if (empty($id)) notice(920201);

        $params['params']['enterprise_id'] = $this->enterpriseId;
        $params['params']['label_id']      = $id;

        return (new GoodsLabel())->deleteGoodsLabel($params);
    }
}
