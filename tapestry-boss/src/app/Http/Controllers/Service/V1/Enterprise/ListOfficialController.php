<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Models\Enterprise;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class ListOfficialController extends BaseController
{
    /**
     * 获取所有正式的企业信息
     * @param Request $request
     * @return array
     * @author: DLong
     */
    public function index(Request $request)
    {

        $enterprise = Enterprise::query()
            ->where([
                ['corp_id', '<>', ''],
                ['status', '=', 1]
            ])
            ->get();

        return success(
            ['data' => $enterprise]
        );
    }
}
