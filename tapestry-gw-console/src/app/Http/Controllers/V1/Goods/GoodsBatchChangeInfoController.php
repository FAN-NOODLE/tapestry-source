<?php


namespace App\Http\Controllers\V1\Goods;

use App\Jobs\Export\ExportGoodsRawJob;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\JsonResponse;

class GoodsBatchChangeInfoController extends BaseController
{
    /**
     * 下载中心route_url
     */
    public const ROUTE_URL = 'goods/goods-batch-change-info-download';

    /**
     * 商品列表出tag
     */
    public const TAG = 'export-goods-raw';

    /**
     * @return JsonResponse
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(): JsonResponse
    {
        $input = $this->getParams();
        $params = [
            'enterprise_hash' => $this->request->header('enterprise-hash'),
            'enterprise_id' => $this->enterpriseId,
            'operate_name'  => $this->enterpriseAdmin['username'] ?? '未知',
            'operate_id'    => $this->enterpriseAdminId ?? 0,
            'route_url'     => self::ROUTE_URL,
            'tag'           => self::TAG,
            'status'        => $input['status'] ?? 1,
        ];
        $this->dispatch(new ExportGoodsRawJob($params, $params, '商品信息'));
        return success();
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'status'    => 'integer',
            'route_url' => 'string',
            'tag'       => 'string',
        ], [
            'status.*'      => '状态不合法',
            'route_url.*'   => '导出路由不合法',
            'tag.*'         => '导出tag不合法',
        ]);
    }

}