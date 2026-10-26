<?php
/**
 * 物料新增
 * Created by PhpStorm.
 * User: 李坚
 * Date: 2020/3/12
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\MaterialRequest;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class AddController extends BaseController
{
    /**
     * 调用接口类型 1-admin 2-api
     */
    const OPERATE_TYPE = 1;

    /**
     * 新增物料
     * @api v1/material/add
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $params                  = $this->request->only([
            'spu_no', // 货品编号
            'material_no', // 物料编号
            'name', // 品名
            'source_id', // 来源id
            'cost_price', // 成本价
            'stock', // 库存
            'weight', // 权重
            'bar_code', // 条码
            'is_channel', // 是否渠道使用 1-是 2-否
            'batch_no', // 批次号
            'desc', // 批次备注
            'taxab_code', // 税目编码，
            'sign_price', // 标牌价
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operate_id']    = $this->enterpriseAdminId;
        $params['operate_name']  = $this->enterpriseAdmin['username'];
        $params['type']          = self::OPERATE_TYPE;

        MaterialRequest::add($params);

        try {
            (new Material())->add($params);
        } catch (Error $error) {
            notice(451512, [], $error->data);
        } catch (\Exception $exception) {
            Log::info('物料库存新增失败' . $exception->getMessage(), $params);
            notice(451512);
        }

        return success();
    }

}
