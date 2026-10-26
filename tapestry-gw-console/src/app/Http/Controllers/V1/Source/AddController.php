<?php
/**
 * 素材添加
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/23
 * Time: 11:04
 */

namespace App\Http\Controllers\V1\Source;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Source;

class AddController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        return (new Source())->sourceAdd($input);
    }
}
