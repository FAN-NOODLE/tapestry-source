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

class CategoryRecursionListController extends BaseController
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
            'cate_tpl_id'
        ]);

        // 参数验证规则
        $rules = [
            'cate_tpl_id' => 'required|integer'
        ];
        //提示
        $messages = [
            'cate_tpl_id.*' => '分类模版id必传',
        ];

        $this->validate($request, $rules, $messages);
        // 获取规格模板下规格信息
        $list = (new CategoryTemplate())->getCategoryRecursionListByTplId($input);

        return success($list);
    }
}