<?php
/**
 * 素材详情
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/23
 * Time: 10:00
 */

namespace App\Http\Controllers\V1\Source;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Source;

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

        return (new Source())->sourceDetail($input);
    }
}
