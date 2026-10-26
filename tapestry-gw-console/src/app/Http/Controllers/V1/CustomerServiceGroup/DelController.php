<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2020-07-28
 * Time: 11:42
 */

namespace App\Http\Controllers\V1\CustomerServiceGroup;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CustomerServiceGroup;
use Illuminate\Http\Request;


class DelController extends BaseController
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

        $this->writeOperatorLog("删除客服分组");

        return (new CustomerServiceGroup())->del($input);
    }
}
