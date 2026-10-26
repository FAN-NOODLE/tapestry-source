<?php
/**
 * 删除我的模板
 * Created by PhpStorm.
 * User: 朱虹达
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\MyTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'id'
        ]);

        // 企业id
        $input['enterprise_id'] = $this->enterpriseId;

        return (new Template())->myDelete($input);
    }
}
