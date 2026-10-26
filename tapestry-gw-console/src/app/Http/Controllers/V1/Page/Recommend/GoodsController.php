<?php
/**
 * Created by PhpStorm.
 * User: yingcai
 * Date: 2019/5/9
 * Time: 9:58 AM
 */

namespace App\Http\Controllers\V1\Page\Recommend;

use App\Common\Constant;
use CrsSdk\Shop\Goods;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class GoodsController extends BaseController
{
    /**
     * 推荐商品选择列表
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author houyingcai
     * @date 2019/05/09 09:59
     */
    public function index(Request $request)
    {
        $params = $request->input();

        // 查询参数处理
        isset($params['spu_name']) && $where['spu_name'] = $params['spu_name'];
        isset($params['no_spu_id']) && $where['no_spu_id'] = $params['no_spu_id'];
        $where['page']          = isset($params['page']) ? $params['page'] : Constant::PAGE;
        $where['limit']         = isset($params['limit']) ? $params['limit'] : Constant::LIMIT;
        $where['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Goods())->activity($where);
            return success($list);
        } catch (\Exception $e) {
            error(710004);
        }
    }
}