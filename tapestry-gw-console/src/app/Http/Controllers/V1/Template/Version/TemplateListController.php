<?php
/**
 * 摸板选择列表数据
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/05/30
 * Time: 11:14 AM
 */

namespace App\Http\Controllers\V1\Template\Version;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use App\Common\Constant;
use CrsSdk\Shop\Version;

class TemplateListController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input['type'] = $request->get('type');
        $input['page'] = $request->get('page', Constant::DEFAULT_PAGE);
        $input['size'] = $request->get('size', Constant::DEFAULT_SIZE);


        // 参数验证规则
        $rules = [
            'type' => 'required'
        ];

        // 提示语
        $message = [
            'type.*' => '类型有误'
        ];

        $this->validate($request, $rules, $message);
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("查看摸板列表数据");

        $list = (new Version())->templateList($input);
        return success($list);
    }
}
