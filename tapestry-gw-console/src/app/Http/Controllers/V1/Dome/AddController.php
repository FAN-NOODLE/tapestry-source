<?php
/**
 * 穹顶添加
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/30
 * Time: 15:00
 */

namespace App\Http\Controllers\V1\Dome;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Dome;

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

        return (new Dome())->domeAdd($input);
    }
}
