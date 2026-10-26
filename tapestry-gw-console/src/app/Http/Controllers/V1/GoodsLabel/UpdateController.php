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

class UpdateController extends BaseController
{

    /**
     * 编辑商品标签
     * {domain}/app/v1/goods-label/update
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 标签id
        $id = request()->post('id', 0);

        // 标签id不能为空
        if (empty($id)) {
            notice(920201);
        }

        // 标签名称
        $name = request()->post('name', '');
        $name = trim($name);
        // 标签名称不能为空
        if (empty($name)) notice(920200);

        // 组装参数
        $params = [
            'label_id'      => $id,
            'name'          => $name,
            'enterprise_id' => $this->enterpriseId
        ];

        return (new GoodsLabel())->updateGoodsLabel($params);
    }
}
