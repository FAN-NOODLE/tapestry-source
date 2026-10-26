<?php
/**
 * Created by PhpStorm.
 * User: zy
 * Date: 2020/3/5
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Category;
use Illuminate\Http\Request;

class StoreCategoryListController extends BaseController
{
    /**
     * 门店
     */
    CONST STORE_TYPE = 2;
    /**
     * 存在SKU的分类
     */
    CONST EXISTS_SKU = 1;

    /**
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index(Request $request)
    {

        // 接受参数
        $input = $request->only([
            'source_id'
        ]);

        // 参数验证规则
        $rules = [
            'source_id' => 'required|integer'
        ];
        //提示
        $messages = [
            'source_id.*' => '门店id必传',
        ];

        $this->validate($request, $rules, $messages);

        //门店
        $input['type'] = self::STORE_TYPE;
        //不展示空的分类
        $input['show_empty'] = 0;
        //展示有sku的分类
        $input['sku_type'] = self::EXISTS_SKU;

        $list = (new Category())->recursionList($input);

        return success($list);
    }
}