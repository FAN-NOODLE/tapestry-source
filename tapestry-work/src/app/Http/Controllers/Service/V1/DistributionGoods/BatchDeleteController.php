<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\DistributionGoods;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use Illuminate\Http\Request;

class BatchDeleteController extends BaseController
{
    /**
     * 分销商品批量删除
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['distribution_goods_ids']);
        $this->validate($request, [
            'distribution_goods_ids' => 'required',
        ], [
            'distribution_goods_ids.required' => '分销商品id不能为空',
        ]);

        DistributionGoods::query()
            ->whereIn('id', $params['distribution_goods_ids'])
            ->delete();

        return success();
    }
}
