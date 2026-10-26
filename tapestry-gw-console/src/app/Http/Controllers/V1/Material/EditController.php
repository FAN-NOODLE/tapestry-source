<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/2
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\MaterialRequest;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class EditController extends BaseController
{
    /**
     * 编辑物料信息
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     *@api v1/material/edit
     */
    public function index()
    {
        $params = $this->request->only([
            'material_id',
            'name',
            'material_no',
            'spu_no',
            'weight',
            'cost_price',
            'sign_price',
            'bar_code',
            'batch_no',
            'desc',
            'source_id',
            'taxab_code',
        ]);

        MaterialRequest::edit($params);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['source_id'] = $params['source_id'] ?? $this->enterpriseId;

        if (empty($params['source_id'])) {
            notice(451521);
        }

        try {
            (new Material())->edit($params);
        } catch (Error $error) {
            Log::info('编辑库存物料失败'.$error->getMessage(), $params);
            throw new Error($error->getCode(), $error->getMessage(), $error->data,210, null, []);
            //throw new Error($error->getCode(), $error->getMessage(), 210, null, [], $error->data);
        } catch (\Exception $exception) {
            Log::info('编辑库存物料失败'.$exception->getMessage(), $params);
            notice(451515);
        }

        return success();
    }

}