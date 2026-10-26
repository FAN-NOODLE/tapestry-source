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

class AddController extends BaseController
{

    /**
     * 新增商品标签
     * {domain}/app/v1/goods-label/add
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 标签名称
        $name = request()->post('name', '');
        $name = trim($name);
        // 标签名称不能为空
        if (empty($name)) notice(920200);

        // 组装参数
        $params['params'] = [
            'name'          => $name,
            'enterprise_id' => $this->enterpriseId,
        ];

        return (new GoodsLabel())->addGoodsLabel($params);
    }
}
