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

class RepeatInspectController extends BaseController
{
    /**
     * 分销商品重复校验
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['sku_ids', 'enterprise_id']);
        $this->validate($request, [
            'sku_ids' => 'required',
            'enterprise_id' => 'required',
        ], [
            'sku_ids.required' => '商品skuid不能为空',
            'enterprise_id.required' => '企业id不能为空',
        ]);

        $list = DistributionGoods::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('sku_id', $params['sku_ids'])
            ->get()
            ->toArray();

        $return = [];
        foreach ($list as $item) {
            $return[$item['sku_id']][] = [
                'distribution_start' => $item['distribution_start'],
                'distribution_end' => $item['distribution_end'],
            ];
        }

        $data = [];
        foreach ($return as $key => $value) {
            array_multisort(array_column($value, 'distribution_start'), SORT_ASC, $value);
            $data[] = [
                'sku_id' => $key,
                'date_range' => $value
            ];
        }

        return success($data);
    }
}
