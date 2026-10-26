<?php

namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use CrsSdk\Shop\Store;

class SourceStockListController extends BaseController
{
    /**
     * 查物料渠道库存列表
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $data = (new MaterialChannel())->sourceStockList($input)['list'];

        // 匹配门店名称
        $storeList = (new Store())->lists([
            'enterprise_id' => $input['enterprise_id'],
            'store_ids'     => array_column($data, 'source_id'),
            'page' => 1,
            'size' => 999999
        ])['list'];

        $storeList = array_column($storeList, 'name', 'id');
        foreach ($data as $key => $value) {
            $data[$key]['store_name'] = $storeList[$value['source_id']] ?? '商城或其它';
        }

        return success([
            'list' => $data
        ]);
    }
}
