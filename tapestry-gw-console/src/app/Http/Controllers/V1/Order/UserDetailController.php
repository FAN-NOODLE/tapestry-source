<?php
/**
 * 拼团用户
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Group;
use Illuminate\Http\Request;

class UserDetailController extends BaseController
{
    /**
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->only(['group_id']);

        return (new Group())->userDetail($input);
    }
}
