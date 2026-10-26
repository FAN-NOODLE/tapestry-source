<?php


namespace App\Http\Controllers\V1\StaffBuyManager;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StaffRight;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AddController extends BaseController
{
    /**
     * 添加允许参加内购的员工及设置
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain']     = $this->qyDomain;

        $this->writeOperatorLog('添加内购员工及设置');

        return (new StaffRight())->add($input);
    }
}
