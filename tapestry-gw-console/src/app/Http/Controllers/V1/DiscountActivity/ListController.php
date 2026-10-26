<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\DiscountActivity;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DiscountActivity;

class ListController extends BaseController
{
    /**
     * 活动列表
     * @author 胡龙阳
     * @date 2020/5/28 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['search', 'page', 'size']);

        $search = $params['search'];
        $search['page'] = $params['page'] ?? 1;
        $search['size'] = $params['size'] ?? 10;
        $search['enterprise_id'] = $this->enterpriseId;

        if (!empty($search['activity_type'])) {
            $search['activity_type'] = [$search['activity_type']];
        }

        // 列表查询
        $list = (new DiscountActivity())->list($search);

        return $list;
    }
}
