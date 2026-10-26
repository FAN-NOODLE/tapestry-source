<?php
/**
 * Created by PhpStorm.
 * User: guokai
 * Date: 2019/04/18
 * Time: 12:02
 */

namespace App\Http\Controllers\V1\Page\ReduceMoney;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageReduce;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 立减金列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 初始化
        $data = [];

        // 接受查询参数
        $data['page'] = $request->post('page', 1);
        $data['size'] = $request->post('size', 15);
        $data['enterprise_id'] = $this->enterpriseId;

        $spuDate = (new PageReduce())->lists($data);

        return success($spuDate);
    }
}
