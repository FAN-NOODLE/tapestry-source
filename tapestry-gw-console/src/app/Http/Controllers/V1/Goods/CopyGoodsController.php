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
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Support\Facades\Log;

class CopyGoodsController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //查询条件校验
        $params = $this->validates();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operator_id'] = $this->enterprise['id'];
        $params['operator_name'] = $this->enterprise['contact'];
        $this->writeOperatorLog('复制商品');
        try {
            (new Goods())->copyGoods($params);
        } catch (Error $error) {

            Log::info('商品创建失败', $error->data);
            notice(450000, [$error->getMessage()]);

        } catch (\Exception $exception) {

            Log::info('商品创建失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }


        //返回结果
        return success();
    }

    /**
     * 数据基础验证
     * @return array
     * @throws \Exception
     */
    private function validates()
    {
        $input = $this->request->all();

        $postageRelationSdk = new PostageRelation();
        $postageList = $postageRelationSdk->relationList([
            'object_id'   => $input['source_id'],
            'object_type' => 2
        ])['list'];
        if (empty($postageList)) notice(450000, ['请完善门店的邮费模板信息']);

        return $input;
    }
}