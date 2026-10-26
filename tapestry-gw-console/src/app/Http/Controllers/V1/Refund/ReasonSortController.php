<?php
/**
 * Created by PhpStorm.
 * User: 售后理由排序
 * Date: 2019/4/4
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Refund;
use Illuminate\Http\Request;

class ReasonSortController extends BaseController
{
    /**
     * 退款原因保存
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $data = $request->input();
        $data['enterprise_id'] = $this->enterpriseId;

        return success((new Refund())->reasonSort($data));
    }
}