<?php


namespace App\Http\Controllers\V1\Supplier;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Supplier;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Supplier())->supplierList($input);
        $list['list'] = $this->formatData($list['list']);

        return success($list);
    }

    private function formatData($list)
    {
        $return = [];
        if (empty($list)) return $return;

        foreach ($list as $item) {
            $return[] = [
                'store_name'   => $item['store_name'],
                'source_id'    => $item['source_id'],
                'warehouse_id' => $item['warehouse_id'],
                'number'       => $item['number'],
            ];
        }

        return $return;
    }
}