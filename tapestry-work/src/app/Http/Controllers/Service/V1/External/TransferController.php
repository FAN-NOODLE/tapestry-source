<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 20:47
 */

namespace App\Http\Controllers\Service\V1\External;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\External;
use Illuminate\Http\Request;

class TransferController extends BaseController
{

    /**
     * 离职成员外部联系人再分配
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain  = $request->post('qy_domain', '');
        $condition = $request->post('condition', []);

        $result = (new External())->transfer($qyDomain, $condition);

        return success($result);
    }
}
