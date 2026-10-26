<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/4
 * Time: 14:08
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Refund;
use Illuminate\Http\Request;

class ReasonListController extends BaseController
{
    /**
     * 退款原因列表
     *
     * @param Request $request
     * @return array
     */
    public function index()
    {
        return success((new Refund())->reasonLists(['enterprise_id' => $this->enterpriseId]));
    }
}
