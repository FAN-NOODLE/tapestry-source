<?php
/**
 * 修改商品是否展示开关
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class UpdateSpuController extends BaseController
{

    /**
     * 商品是否展示开关
     * @api v1/material/add
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接收数据数据验证
        $params = $this->getParams($request);

        $input = [
            'spu_id' => $params['spu_id'],
            'save' => [
                'is_hide' => $params['is_hide']
            ],
        ];
        $this->writeOperatorLog('修改商品的展示状态');
        try {
            (new Goods())->updateSpu($input);
        } catch (\Exception $exception) {
            Log::info('商品spu操作失败'.$exception->getMessage(), $params);
            notice(451502);
        }

        return success();
    }

    /**
     * 参数验证
     *
     * @param Request $request
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'spu_id'  => 'required|integer',
            'is_hide' => 'required|integer',
        ];

        // 提示语
        $messages = [
            'spu_id.*'  => 'spuId有误',
            'is_hide.*' => '开关类型有误',
        ];

        $this->validate($request, $rules, $messages);

        return $input;
    }

}
