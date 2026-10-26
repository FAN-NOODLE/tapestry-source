<?php
/**
 * 修改售后单状态
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/07/11
 * Time: 18:20
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Refund;

class ChangeRefundApplyStateController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operate_store_id'] = $this->enterpriseId; //todo 一期先用企业id
        $input['operate_id'] = $this->enterpriseAdmin['id'];
        $input['operate_name'] = $this->enterpriseAdmin['username'];
        $this->writeOperatorLog("后台改变售后单状态");

        $response = (new Refund())->machineRun($input);
        return success($response);
    }
}