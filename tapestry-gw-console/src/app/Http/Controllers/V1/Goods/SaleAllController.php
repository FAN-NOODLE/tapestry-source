<?php

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use App\Services\PostageService;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use App\Exceptions\Error;
use CrsSdk\MicroGoodsV1\WarehouseStock;
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class SaleAllController extends BaseController
{
    /**
     * 上下架状态 1：下架
     */
    const SALE_DOWN = 1;

    /**
     * 上下架状态 2：上架
     */
    const SALE_UP = 2;

    /**
     * 库存状态 0：充足
     */
    const ALL_ENOUGH = 0;

    /**
     * 库存状态 1：线上库存不足
     */
    const ON_LINE_LESS = 1;

    /**
     * 库存状态 2：线下库存不足
     */
    const OFF_LINE_LESS = 2;

    /**
     * 库存状态 3：总库存不足
     */
    const ALL_LESS = 3;

    /**
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        $input                  = $this->getParams();
        $input['enterpriseId']  = $this->enterpriseId;
        $input['page'] = 1;

        $request = [];
        $request['status'] = 1;
        $request['operate_id'] = $this->enterpriseAdminId ?? 1;
        $request['source_id'] = $input['source_id'] ?? $this->enterpriseId;
        $request['data'] = [];
        do{
            // 获取商品列表
            $list = (new Goods())->getConsoleList($input);
            $arr = $list['list'];
            foreach($arr as $value) {
                $request['data'][] = [
                    'source_id' => $value['source_id'],
                    'spu_id' => $value['spu_id'],
                    'sku_id' => $value['sku_id'],
                    'material_channel_id' => $value['channel_id']
                ];
            }
            $input['page'] += 1;
        }while(count($arr) != 0);

        $ret = $this->saleGoods($request);
        return success([
            'fail_sku'    => $ret['fail_sku'],
            'success_num' => $ret['success_num'],
            'fail_num'    => $ret['fail_num'],
            'fail_goods'  => $ret['fail_goods']
        ]);
    }

    /**
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate(request(), [
            'searchType'    => 'required|integer',
            'searchParams.searchString.keyword' => 'string',
            'searchParams.searchString.type'    => 'integer',
            'searchParams.searchInt.high'       => 'integer',
            'searchParams.searchInt.low'        => 'integer',
            'searchParams.searchInt.type'       => 'integer',
            'searchParams.searchChannel'        => 'integer',
            'searchParams.searchCreateTime.beginTimestamp'  => 'integer',
            'searchParams.searchCreateTime.endTimestamp'    => 'integer',
            'searchParams.searchGoodsStatus'    => 'integer',
            'searchParams.searchChannelOwnership'   => 'integer',
            'searchParams.searchSort.sortColumn'    => 'integer',
            'searchParams.searchSort.sortType'      => 'integer',
            'page'  => 'integer',
            'size'  => 'integer'
        ], [
            'searchType.*'      => '列表类型不合法',
            'searchParams.searchString.keyword.*'   => '关键字必须是字符串',
            'searchParams.searchString.type.*'      => '关键字搜索类型不合法',
            'searchParams.searchInt.high.*'     => '最大值不合法',
            'searchParams.searchInt.low.*'      => '最小值不合法',
            'searchParams.searchInt.type.*'     => '范围类型不合法',
            'searchParams.searchChannel.*'      => '渠道不合法',
            'searchParams.searchCreateTime.beginTimestamp.*'    => '开始时间不合法',
            'searchParams.searchCreateTime.endTimestamp.*'      => '结束时间不合法',
            'searchParams.searchGoodsStatus.*'  => '商品状态不合法',
            'searchParams.searchChannelOwnership.*' => '渠道归属不合法',
            'searchParams.searchSort.sortColumn.*'  => '排序字段不合法',
            'searchParams.searchSort.sortType.*'    => '排序类型不合法',
            'page.*'    => '分页不合法',
            'size.*'    => '分页大小不合法',
        ]);
    }


    public function saleGoods ($input)
    {
        // 记录操作日志
        $this->writeOperatorLog('管理员'.$this->enterpriseAdmin['username'].'操作:商品全部下架');
        // 下架不区分单个还是批量
        $failList = [];
        $message = '';
        $successNum = 0;
        $failNum = 0;
        $materialSkd = new MaterialChannel();

        foreach ($input['data'] as $item) {
            try {
                $detail = $materialSkd->detail([
                    'enterprise_id'       => $this->enterpriseId,
                    'source_id'           => $item['source_id'],
                    'material_channel_id' => $item['material_channel_id'],
                ]);
                $this->flag($item, $input['status'], $detail);
                $successNum++;
            } catch (Error $error) {
                $failNum++;
                Log::error('上下架失败', $error->data);
                $failList[] = $item['sku_id'];
                $message = $error->getMessage();
            } catch (\Exception $e) {
                $failNum++;
                Log::error('上下架失败' . $e->getMessage(), $input);
                $failList[] = $item['sku_id'];
                $message = $e->getMessage();
            }
        }
        return [
            'fail_sku'    => $failList,
            'success_num' => $successNum,
            'fail_num'    => $failNum,
            'fail_goods'  => $message
        ];
    }

    /**
     * 上下架
     *
     * @param $input
     */
    public function flag($item, $status, $detail)
    {
        $materialSkd = new MaterialChannel();

        if (empty($detail)) {
            notice(450000, ['未找到渠道信息']);
        }

        $params = [
            'object_id'     => $item['sku_id'],
            'source_id'     => $item['source_id'],
            'status'        => $status,
            'operate_id'    => $this->enterpriseAdminId ?? 1,
            'enterprise_id' => $this->enterpriseId,
            'data'          => [
                [
                    'sku_id'                 => $item['sku_id'],
                    'on_line_current_stock'  => $detail['on_line_current_stock'],
                    'off_line_current_stock' => $detail['off_line_current_stock'],
                    'sale_type'              => $detail['sale_type']
                ]
            ]
        ];

        $materialSkd->flag($params);

        return true;
    }


    /**
     * 关联商品的邮费模板信息
     * @param $data
     * @return int
     * @throws \Exception
     */
    protected function postageRelation($data)
    {
        $postageRelation = new PostageRelation();

        $thisRelation = $postageRelation->relationList([
            'object_id'     => $data['source_id'],
            'object_type'   => 2,
            'enterprise_id' => $this->enterpriseId
        ])['list'];

        if (empty($thisRelation)) notice(450000, ['请完善门店的邮费模板信息']);

        array_map(function ($v) use (&$skuRelation, $data) {
            $skuRelation[] = [
                'postage_id'    => $v['postage_id'],
                'delivery_type' => $v['delivery_type']
            ];
        }, $thisRelation);
        $postData = [
            'object_type' => 1,
            'object_id'   => $data['sku_id'],
            'relation'    => $skuRelation,
            'enterprise_id' => $this->enterpriseId
        ];
        $postageRelation->postageRelationSkuOrActivity($postData);

        // 删除商品缓存
        (new Goods())->deleteRedis([
            'sku_id'        => $data['sku_id'],
            'enterprise_id' => $this->enterpriseId
        ]);

        $deliveryArr = array_column($skuRelation, 'delivery_type');

        return (new PostageService())->getSaleTypeByDeliver($deliveryArr);
    }
}
