<?php


namespace App\Http\Controllers\V1\Supplier;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Supplier;
use CrsSdk\Shop\Store;

class StoreListController extends BaseController
{
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $this->validate($this->request, [
            'warehouse_id' => 'required|integer|min:1',
        ], [
            'warehouse_id.*' => '仓库ID不能为空'
        ]);

        $storeList = (new Store())->lists([
            'enterprise_id' => $this->enterpriseId,
            'size'          => 999999
        ])['list'];
        // 将商城拼入门店列表
        $storeList = array_merge([
            [
                'id'   => $this->enterpriseId,
                'name' => '商城'
            ]
        ], $storeList);
        $storeList = array_map(function ($element) {
            return [
                'id'   => $element['id'],
                'name' => $element['name'],
            ];
        }, $storeList);

        // 取出现有的供货关系门店
        $input['size'] = 999999;
        $supplierList = (new Supplier())->supplierList($input)['list'];
        $sourceIds = array_column($supplierList, 'source_id');

        foreach ($storeList as $key => $item) {
            if (in_array($item['id'], $sourceIds)) {
                unset($storeList[$key]);
            }
        }

        return success([
            'list'  => array_values($storeList),
            'total' => count($storeList),
            'page'  => 1,
            'size'  => 99999
        ]);
    }

}