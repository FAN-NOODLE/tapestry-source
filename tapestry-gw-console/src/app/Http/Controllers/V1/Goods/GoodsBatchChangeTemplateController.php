<?php


namespace App\Http\Controllers\V1\Goods;


use App\Http\Controllers\V1\BaseController;
use App\Services\GoodsBatchChangeService;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

class GoodsBatchChangeTemplateController extends BaseController
{
    /**
     * @return BinaryFileResponse
     * @throws ValidationException
     */
    public function index(): BinaryFileResponse
    {
        $path = '';
        $params = $this->getParams();
        switch ($params['change_type']) {
            case GoodsBatchChangeService::TYPE_WAIT_SALES_PRICE:
                $path = resource_path('static/待上架商品更改价格模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_GOODS_NAME_SPU:
                $path = resource_path('static/SPU纬度修改商品名称模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_GOODS_NAME_SKU:
                $path = resource_path('static/SKU纬度修改商品名称模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_WAIT_SALES_STOCK_EXPRESS:
                $path = resource_path('static/待上架商品修改线上库存模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_WAIT_SALES_STOCK_EXTRACT:
                $path = resource_path('static/待上架商品修改线下库存模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_ALREADY_SALES_STOCK_EXPRESS:
                $path = resource_path('static/上架商品修改线上库存模板.xlsx');
                break;
            case GoodsBatchChangeService::TYPE_ALREADY_SALES_STOCK_EXTRACT:
                $path = resource_path('static/上架商品修改线下库存模板.xlsx');
                break;
        }
        return response()->download($path);
    }

    /**
     * @return array
     * @throws ValidationException
     */
    protected function getParams(): array
    {
        $inStr = GoodsBatchChangeService::TYPE_COLLECT_STR;
        return $this->validate(request(), [
            'change_type'   => 'required|integer|in:' . $inStr,
        ], [
            'change_type.*' => '更改类型不合法',
        ]);
    }
}