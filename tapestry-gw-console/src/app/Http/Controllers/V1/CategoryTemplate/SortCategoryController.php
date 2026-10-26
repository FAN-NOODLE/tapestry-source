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
use Illuminate\Http\Request;

class SortCategoryController extends BaseController
{


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
            'category_id',
            'type'
        ]);

        // 参数验证规则
        $rules = [
            'category_id' => 'required|integer',
            'type'        => 'required|in:1,2'
        ];


        $this->validate($request, $rules);

        // 分类移动
        $res = (new CategoryTemplate)->sortCategory($input);
        return success($res);

    }
}