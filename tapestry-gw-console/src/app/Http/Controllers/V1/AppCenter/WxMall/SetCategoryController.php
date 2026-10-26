<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019-07-16
 * Time: 14:00
 */

namespace App\Http\Controllers\V1\AppCenter\WxMall;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxMall;
use Illuminate\Http\Request;

class SetCategoryController extends BaseController
{
    /**
     * 好物圈类目列表接口
     * {domain}/service/v1/wx-mall/category-list
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params                  = $this->getParams($request);
        $params['enterprise_id'] = $this->enterpriseId;

        (new WxMall())->wxMallSetCategory($params);

        return success(['data' => []]);
    }

    protected function getParams(Request $request)
    {
        // 接受参数
        $request->input(['spu_ids', 'category_list']);

        // 参数验证规则
        $rules = [
            'category_id' => 'required',
            'spu_ids'       => 'required|array',
        ];

        // 提示语
        $message = [
            'category_id.*' => '好物圈类目必填',
            'spu_ids.*'       => '商品spu_ids必填，数组格式'
        ];

        $input = $this->validate($request, $rules, $message);

        return $input;
    }
}
