<?php
/**
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

class DetailController extends BaseController
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
            'spu_ids'       => [$param['spu_id']]
        ];

        try {
            $spuInfo = (new Goods())->spuDetail($input);
        } catch (Error $error) {

            Log::info('获取商品详情失败', [$error->getMessage()]);
            notice(451503, [], [$error->getMessage()]);

        } catch (\Exception $exception) {

            Log::info('获取商品详情失败' . $exception->getMessage(), $input);
            notice(451503);
        }

        //返回结果
        return success($spuInfo);
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
            'spu_id' => 'required|integer',
        ];

        $message = [
            'condition.spu_id.*' => '商品ID有误',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }
}