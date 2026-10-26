<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2020/2/20                   
* +------------------------------------+
* | Time   : 16:39                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\V1\DoveDream;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Validator;
use Illuminate\Http\Request;
use CrsSdk\Shop\DoveDream;


class SkuController extends BaseController
{

    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 商品库存(sku)列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function list(Request $request)
    {
        $params = $request->only(['spu_id', 'year', 'month']);

        $rules = [
            'spu_id' => 'required',
            'year'   => 'required|integer',
            'month'  => 'required|integer',
        ];

        $message = [
            'spu_id.required' => '商品id不能为空',
            'year.required'   => '库存年数不能为空',
            'month.required'  => '库存月数不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getSkuList($params);

        return success($result);
    }


    /**
     * 商品库存(sku)保存
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function save(Request $request)
    {
        $params = $request->only(['spu_id', 'year', 'month', 'is_batch', 'stocks']);

        $rules = [
            'spu_id'               => 'required',
            'year'                 => 'required|integer',
            'month'                => 'required|integer',
            'stocks'               => 'required|array',
            'stocks.*.date'        => 'required|date',
            'stocks.*.am_sku_id'   => 'required|min:0',
            'stocks.*.pm_sku_id'   => 'required|min:0',
            'stocks.*.am_is_full'  => 'required|integer|in:0,1',
            'stocks.*.pm_is_full'  => 'required|integer|in:0,1',
            'stocks.*.am_is_close' => 'required|integer|in:0,1',
            'stocks.*.pm_is_close' => 'required|integer|in:0,1',
        ];

        $message = [
            'spu_id.required'               => '商品id不能为空',
            'year.required'                 => '库存年数不能为空',
            'month.required'                => '库存月数不能为空',
            'stocks.required'               => '库存数据不能为空',
            'stocks.array'                  => '库存数据类型错误',
            'stocks.*.date.required'        => '每日日期不能为空',
            'stocks.*.am_sku_id.required'   => '上午库存id不为空',
            'stocks.*.pm_sku_id.required'   => '下午库存id不为空',
            //            'stocks.*.am_stock.required'    => '上午库存数不为空',
            //            'stocks.*.pm_stock.required'    => '下午库存数不为空',
            'stocks.*.am_is_full.required'  => '上午约满状态不为空',
            'stocks.*.pm_is_full.required'  => '下午约满状态不为空',
            'stocks.*.am_is_close.required' => '上午闭馆状态不为空',
            'stocks.*.pm_is_close.required' => '下午闭馆状态不为空',


        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->saveSku($params);

        return success($result);
    }
}