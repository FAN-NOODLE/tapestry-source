<?php
/**
 * 版本记录列表
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/06/13
 * Time: 16:14
 */

namespace App\Http\Controllers\V1\Template\Version;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Version;

class SaveStatusController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input['status']        = '2';
        $input['id']            = $request->get('id');

        return (new Version())->saveStatus($input);

    }
}
