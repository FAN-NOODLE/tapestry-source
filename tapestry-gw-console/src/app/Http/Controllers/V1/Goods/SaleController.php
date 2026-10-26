<?php

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use App\Services\PostageService;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use App\Exceptions\Error;
use CrsSdk\MicroGoodsV1\WarehouseStock;
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Support\Facades\Log;

class SaleController extends BaseController
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
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->validateInput();

        // 处理spu上架
        $this->checkSpuSale($input);

        // 区分单个上架还是批量上架
        if (self::SALE_UP == $input['status']) {
            $this->writeOperatorLog('上架商品');
            if (1 == count($input['data'])) {
                $return = $this->oneFlag($input);
            } else {
                $return = $this->manyFlag($input);
            }
        } else {
            $this->writeOperatorLog('下架商品');
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

            $return = [
                'fail_sku'    => $failList,
                'success_num' => $successNum,
                'fail_num'    => $failNum,
                'fail_goods'  => $message
            ];

        }

        return success($return);
    }

    /**
     * @return array|null|string
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function validateInput()
    {
        $input = $this->request->input();

        // 验证数据
        $rule = [
            'status'                     => "required|integer|in:" . self::SALE_UP . ',' . self::SALE_DOWN,
            'data'                       => 'required|array',
            'data.*.source_id'           => 'required|integer',
            'data.*.spu_id'              => 'required_without:data.*.sku_id',
            'data.*.sku_id'              => 'required_without:data.*.spu_id',
            'data.*.material_channel_id' => 'required_with:data.*.sku_id',
            'data.*.name'                => 'string',
        ];

        $message = [
            'status.required'              => '上下架状态不能为空',
            'status.integer'               => '上下架状态必须是数字',
            'status.in'                    => '上下架状态(1,下架 2,上架)',
            'data.*.source_id.*'           => '来源ID不能为空',
            'data.*.spu_id.*'              => '请选择商品',
            'data.*.sku_id.*'              => '请选择商品',
            'data.*.material_channel_id.*' => '请选择渠道',
        ];

        $this->validate($this->request, $rule, $message);

        $input['operate_id'] = $this->enterpriseAdminId ?? 1;
        $input['source_id'] = $input['source_id'] ?? $this->enterpriseId;

        return $input;
    }

    /**
     * @param $input
     * @throws \Exception
     */
    public function checkSpuSale(&$input)
    {
        // 如果是spu上架，获取spu下的多个sku (spu无批量上下架)
        $tempSpu = [];
        foreach ($input['data'] as $key => $value) {
            if (isset($value['spu_id']) && !empty($value['spu_id'])) {
                $tempSpu[] = $value;
                unset($input['data'][$key]);
            }
        }

        $skuList = !empty($tempSpu) ? (new Goods())->spuDetail([
            'spu_ids'       => array_column($tempSpu, 'spu_id'),
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => $input['source_id']
        ])['list'] : [];

        if (isset($skuList[0]['sku']) && !empty($skuList[0]['sku'])) {
            foreach ($skuList[0]['sku'] as $sku) {
                if ($sku['status'] != $input['status']) {
                    $input['data'][] = [
                        'sku_id'              => $sku['id'],
                        'source_id'           => $sku['source_id'],
                        'material_channel_id' => $sku['material_channel_id'],
                    ];
                }
            }
        }
    }

    /**
     * 单个上架商品验证库存
     *
     * @param $detail
     *
     * @return int
     */
    protected function checkStockAndPrice($detail)
    {
        // 查询sku线上线下售卖方式
        $saleType = $detail['sale_type'];

        // 查询物料库总库存
        $material = new WarehouseStock();
//        $materialStockData = $material->detail(['material_id' => $detail['material_id']]);
        $materialStockData = $material->getStock([
            'material_id'   => $detail['material_id'],
            'enterprise_id' => $detail['enterprise_id'],
            'source_id'     => $detail['source_id']
        ]);

        if (empty($materialStockData)) notice(450000, ['未找到物料信息']);

        $onlineStock  = $materialStockData['on_line_stock']  + $materialStockData['public_stock'];
        $offlineStock = $materialStockData['off_line_stock'] + $materialStockData['public_stock'];

        // 返回单个上架库存验证状态
        $failType = self::ALL_ENOUGH;
        switch ($saleType) {
            // 仅线上售卖
            case PostageService::ON_LINE_SALE:
                // 线上价格必须大于0
                if ($detail['on_line_price'] <= 0) {
                    notice(450000, ['线上价格必须大于0']);
                }
                if ($onlineStock < $detail['on_line_current_stock']) {
                    $failType = self::ON_LINE_LESS;
                }
                break;
            // 仅线下自提
            case PostageService::OFF_LINE_SALE:
                //线下价格必须大于0
                if ($detail['off_line_price'] <= 0) {
                    notice(450000, ['线下价格必须大于0']);
                }
                if ($offlineStock < $detail['off_line_current_stock']) {
                    $failType = self::OFF_LINE_LESS;
                }
                break;
            // 线上线下均有售
            case PostageService::ALL_SALE:
                // 线上线下价格不能都<=0
                if ($detail['on_line_price'] <= 0 && $detail['off_line_price'] <= 0) {
                    notice(450000, ['线上或线下价格有误']);
                }

                if (($detail['on_line_price'] > 0 && $onlineStock  < $detail['on_line_current_stock'])
                || ($detail['off_line_price'] > 0 && $offlineStock < $detail['off_line_current_stock'])) {
                    $failType = self::ALL_LESS;
                }

                break;
        }

        return $failType;
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
     * 单个上架
     *
     * @param $input
     *
     * @throws \Exception
     */
    public function oneFlag($input)
    {
        $materialSkd = new MaterialChannel();

        $data = current($input['data']);

        $detail = $materialSkd->detail([
            'enterprise_id'       => $this->enterpriseId,
            'source_id'           => $data['source_id'],
            'material_channel_id' => $data['material_channel_id']
        ]);

        if (empty($detail)) notice(450000, ['未找到物料信息']);

        // 不是商城的商品，需要将门店的邮费模板更新到商品
        if ($data['source_id'] != $this->enterpriseId) {
            $detail['sale_type'] = $this->postageRelation($data);
        }

        $failType = $this->checkStockAndPrice($detail);

        if (self::ALL_ENOUGH !== $failType) {
            return [
                'fail_sku'    => [$detail['sku_id']],
                'success_num' => 0,
                'fail_num'    => 1,
                'fail_goods'  => '商品库存不足',
                'fail_type'   => $failType
            ];
        }

        try {
            $this->flag($data, $input['status'], $detail);
        } catch (\Exception $e) {
            Log::error('上下架失败' . $e->getMessage(), $input);
            notice(450000, [$e->getMessage()]);
        }

        try {
            // 普通商品上架成功，如果是第一次上架需要判断是否需要发送上新通知订阅消息
            if ($detail['status'] == 1 && $detail['channel_type'] == 1) {
                $materialSkd->flagSendMessage([
                    'enterprise_id' => $this->enterpriseId,
                    'sku_ids'       => [
                        $data['sku_id']
                    ]
                ]);
            }
        } catch (\Exception $exception) {
            Log::error('发送上新订阅消息失败：' . $exception->getMessage(), [
                'enterprise_id' => $this->enterpriseId,
                'sku_ids'       => [$data['sku_id']]
            ]);
        }

        return [
            'fail_sku'    => [],
            'success_num' => 1,
            'fail_num'    => 0,
            'fail_goods'  => ''
        ];
    }

    /**
     * 批量上下架
     *
     * @param $input
     */
    public function manyFlag($input)
    {
        $failList = [];
        $message = '';
        $successNum = 0;
        $failNum = 0;
        $normalGoodsIds = [];
        $materialSkd = new MaterialChannel();

        foreach ($input['data'] as $item) {
            try {
                $detail = $materialSkd->detail([
                    'enterprise_id'       => $this->enterpriseId,
                    'source_id'           => $item['source_id'],
                    'material_channel_id' => $item['material_channel_id']
                ]);

                if (empty($detail)) notice(450000, ['未找到物料信息']);
                // 不是商城的商品，需要将门店的邮费模板更新到商品
                if ($item['source_id'] != $this->enterpriseId) {
                    $detail['sale_type'] = $this->postageRelation($item);
                }
                // 检查商品的线上线下价格
                $this->checkStockAndPrice($detail);
                // 上架
                $this->flag($item, $input['status'], $detail);
                if ($detail['status'] == 1 && $detail['channel_type'] == 1) {
                    $normalGoodsIds[] = $detail['sku_id'];
                }
                $successNum++;
            } catch (\Exception $e) {
                $failNum++;
                Log::error('上下架失败' . $e->getMessage(), $input);
                $failList[] = $item['sku_id'];
                $message = '上下架失败，请检查商品的信息';
            }
        }

        try {
            if (!empty($normalGoodsIds)) {
                $materialSkd->flagSendMessage([
                    'enterprise_id' => $this->enterpriseId,
                    'sku_ids'       => $normalGoodsIds
                ]);
            }
        } catch (\Exception $exception) {
            Log::error('发送上新订阅消息失败：' . $exception->getMessage(), [
                'enterprise_id' => $this->enterpriseId,
                'sku_ids'       => $normalGoodsIds
            ]);
        }

        return [
            'fail_sku'    => $failList,
            'success_num' => $successNum,
            'fail_num'    => $failNum,
            'fail_goods'  => $message
        ];
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
