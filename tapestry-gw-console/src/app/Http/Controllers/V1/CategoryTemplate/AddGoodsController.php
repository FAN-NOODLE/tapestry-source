<?php
/**
 * Created by PhpStorm.
 * User: zy
 * Date: 2020/3/5
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CategoryTemplate;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class AddGoodsController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return array
     * @throws ValidationException
     * @throws Exception
     */
    public function index(Request $request)
    {

        // 接受参数
        $input = $request->only([
            'category_id',
            'goods_ids'
        ]);

        // 参数验证规则
        $rules = [
            'category_id' => 'required|integer',
            'goods_ids'       => 'required|array'
        ];

        $this->validate($request, $rules);


        // 分类新增商品
        $res = (new CategoryTemplate())->addGoods($input);

        return success($res);
    }
}