<?php
/**
 * 删除商品sku
 *
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class DelSkuController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //查询条件校验
        $param = $this->validates();
        $input = [
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => $this->enterpriseId,
            'channel_type'  => $param['channel_type'],
            'sku_ids'       => $param['sku_ids']
        ];
        $this->writeOperatorLog('删除sku');
        try {
            (new Goods())->deleteSku($input);

        } catch (Error $error) {

            Log::info('删除sku失败', $error->data);
            notice(451505, [], $error->data);

        } catch (\Exception $exception) {

            Log::info('删除sku失败' . $exception->getMessage(), $input);
            notice(451505);
        }

        //返回结果
        return success([]);
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'channel_type' => 'required|integer',
            'sku_ids'      => 'required|array',
        ];

        $message = [
            'channel_type.*' => '商品渠道类型有误',
            'sku_ids.*'      => '商品ID集合有误',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }
}