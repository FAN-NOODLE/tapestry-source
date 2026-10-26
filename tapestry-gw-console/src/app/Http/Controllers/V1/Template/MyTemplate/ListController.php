<?php
/**
 * 我的模板列表
 * Created by PhpStorm.
 * User: 朱虹达
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\MyTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'page',
            'size'
        ]);

        $input['size'] = $input['size'] ?? 10;
        $input['enterprise_id'] = $this->enterpriseId;

        // 查询列表
        return (new Template())->myList($input);
    }
}
