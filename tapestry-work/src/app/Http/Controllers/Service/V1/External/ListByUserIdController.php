<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 20:46
 */

namespace App\Http\Controllers\Service\V1\External;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\External;
use Illuminate\Http\Request;

class ListByUserIdController extends BaseController
{

    /**
     * 根据userid获取外部联系人列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain  = $request->input('qy_domain', '');
        $condition = $request->input('condition', []);

        $result = (new External())->listByUserId($qyDomain, $condition);

        return success($result);
    }
}
