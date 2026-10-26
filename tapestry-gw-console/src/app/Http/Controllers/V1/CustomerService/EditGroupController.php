<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2020-07-28
 * Time: 11:42
 */

namespace App\Http\Controllers\V1\CustomerService;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CustomerService;
use Illuminate\Http\Request;

/**
 * Class StatusController
 *
 * @package App\Http\Controllers\V1\CustomerService
 */
class EditGroupController extends BaseController
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

        $this->writeOperatorLog("编辑客服分组");

        return (new CustomerService())->editGroup($input);
    }
}
