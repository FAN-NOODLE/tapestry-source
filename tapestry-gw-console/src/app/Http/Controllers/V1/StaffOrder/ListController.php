<?php
/**
 * 内购订单列表
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\StaffOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffOrder;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    /**
     * 订单列表接口
     * {domain}/app/v1/order/list
     *
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     */
    public function index(Request $request)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['search']['qy_domain'] = $this->qyDomain;
        $input['search']['app_identifier'] = $this->appIdentifier;

        $data= (new StaffOrder())->lists($input);

        return  $data;
    }
}
