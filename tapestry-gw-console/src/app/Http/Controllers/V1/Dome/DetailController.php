<?php
/**
 * 穹顶详情
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/30
 * Time: 15:00
 */

namespace App\Http\Controllers\V1\Dome;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Dome;

class DetailController extends BaseController
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

        return success((new Dome())->domeDetail($input));
    }
}
