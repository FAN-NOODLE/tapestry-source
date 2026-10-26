<?php
/**
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2020/3/20
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class SpuDraftListController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //查询条件校验
        $input = $this->request->all();
        $input['enterpriseId'] = $this->enterpriseId;
        $input['source_id'] = $this->enterpriseId;

        try {
            $list = (new Goods())->spuDraftList($input);
        } catch (Error $error) {
            Log::info('获取商品草稿列表失败', $error->data);
            notice(451509, [], $error->data);

        } catch (\Exception $exception) {
            Log::info('获取商品草稿列表失败' . $exception->getMessage(), $input);
            notice(451509);
        }
        //返回结果
        return success($list);
    }
}