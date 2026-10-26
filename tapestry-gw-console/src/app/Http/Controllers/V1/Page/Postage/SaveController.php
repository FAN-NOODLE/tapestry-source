<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/3
 * Time: 16:37
 */

namespace App\Http\Controllers\V1\Page\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GlobalPostage;
use Illuminate\Http\Request;

class SaveController extends BaseController
{
    /**
     * 邮费管理保存
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $data['on_off'] = $request->input('on_off');
        $data['rule']   = $request->input('rule');
        $data['enterprise_id'] = $this->enterpriseId;

        (new GlobalPostage())->save($data);

        return success();
    }
}
