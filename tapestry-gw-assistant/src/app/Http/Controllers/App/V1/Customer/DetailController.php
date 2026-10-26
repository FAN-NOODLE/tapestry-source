<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/17
 * Time: 17:59
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class DetailController extends BaseController {
    /**
     * index
     * @desc 顾客详情信息，包含订单数据
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @date 2020-07-20
     * @throws \Exception
     */
    public function index(Request $request) {
        $param['union_id'] = $request->get('union_id');
        $param['store_id'] = $request->get('store_id');
        $param['relation_id'] = $request->get('relation_id');
        $param['mobile'] = $this->loginMobile();
        $param['enterprise_id'] = $this->enterpriseId;
        $param['page'] = $request->get('page', Constant::PAGE);
        $param['size'] = $request->get('size', Constant::LIMIT);
        $result = (new Customer())->customerDetailWithOrder($param);
        return success($result);
    }


    /**
     * 顾客详情(新)
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();

        $result = (new Customer())->customerDetailWithOrder($params);

        return success($result);
    }

}