<?php
/**
 * 删除商品spu
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

class DelSpuController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //查询条件校验
        $param = $this->validates();
        $input = [
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => $this->enterpriseId,
            'spu_id'        => $param['spu_id'],
            'channel_type'  => $param['channel_type']
        ];
        $this->writeOperatorLog('删除spu');
        try {
            (new Goods())->deleteSpu($input);

        } catch (Error $error) {

            Log::info('删除spu失败', $error->data);
            notice(451504, [], $error->data);
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
            'spu_id'       => 'required|integer',
            'channel_type' => 'required|integer',
        ];

        $message = [
            'spu_id.*'       => '商品ID有误',
            'channel_type.*' => '渠道信息类型有误',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }
}