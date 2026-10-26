<?php
/**
 * Created by PhpStorm.
 * Function: xxx
 * User: 龚国新
 * Date: 2020/3/5
 * Time: 15:37
 */

namespace App\Http\Controllers\V1\Label\Goods;

use App\Http\Controllers\V1\Label\LabelBaseController;
use CrsSdk\MicroGoodsV1\Label;

class ListController extends  LabelBaseController
{
    /**
     * 标签列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new Label())->list($params);

        //返回结果
        return success($result);
    }

}