<?php
/**
 * 一键转人工
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2020年2月5日
 * Time: 10:25:30
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Refund;

class RefundToManualController extends BaseController
{
    public function index()
    {
        $input = $this->request->all();
        $this->writeOperatorLog("后台一键转人工");
        $res = (new Refund())->RefundToManual($input);

        return success($res);
    }
}