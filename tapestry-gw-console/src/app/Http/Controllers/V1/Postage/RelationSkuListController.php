<?php


namespace App\Http\Controllers\V1\Postage;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\PostageRelation;
use CrsSdk\Shop\Store;

class RelationSkuListController extends BaseController
{
    /**
     * 搜索类型 1-按照sku_id搜索 2-按照商品名称搜索 3-按照商品spu_id搜索
     */
    const SEARCH_TYPE_SKU = 1;
    const SEARCH_TYPE_NAME = 2;
    const SEARCH_TYPE_SPU = 3;

    const MAX_SIZE = 999;

    /**
     * 获取关联商品列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'postage_id',
            'search_type',
            'keyword',
            'source_id',
            'page',
            'size',
        ]);

        // 搜索条件构建
        $input['page'] = $input['page'] ?? Constant::PAGE;
        $input['size'] = $input['size'] ?? Constant::LIMIT;
        $input['enterprise_id'] = $this->enterpriseId;
        if (isset($input['search_type']) && isset($input['keyword']) && trim($input['keyword']) !== '') {
            switch ($input['search_type']) {
                case self::SEARCH_TYPE_NAME:
                    // 商品名称模糊搜索
                    $input['sku_name'] = trim($input['keyword']);
                    break;
                case self::SEARCH_TYPE_SKU:
                    // 商品skuID搜索
                    $input['sku_id'] = trim($input['keyword']);
                    break;
                case self::SEARCH_TYPE_SPU:
                    // 商品spuID搜索
                    $input['spu_id'] = trim($input['keyword']);
                    break;
            }
        }

        // 获取商品信息
        $data = (new PostageRelation())->postageRelationSkuList($input);

        // 获取souce_id对应的门店
        $sourceIds   = array_values(array_unique(array_column($data['list'], 'source_id')));
        $storeParams = [
            'store_ids'     => $sourceIds,
            'enterprise_id' => $input['enterprise_id'],
            'page'          => Constant::PAGE,
            'size'          => self::MAX_SIZE
        ];
        $storeList   = (new Store())->lists($storeParams);
        $storeList   = array_column($storeList['list'], null, 'id');

        // 处理数据
        $result = [];
        foreach ($data['list'] as $key => $val) {
            $result[] = [
                'postage_id' => $input['postage_id'],
                'spu_id'     => $val['spu_id'],
                'sku_id'     => $val['sku_id'],
                'sku_name'   => $val['sku_name'] ?? '',
                'store_name' => isset($storeList[$val['source_id']]) ? $storeList[$val['source_id']]['name'] : '商城或其它',
                'status'     => $val['status'],
                'sku_image'  => $val['sku_image'],
                'attributes' => $val['attributes']
            ];
        }

        return success([
            'list'  => $result ?? [],
            'page'  => $input['page'],
            'size'  => $input['size'],
            'total' => $data['total'] ?? 0
        ]);
    }

}