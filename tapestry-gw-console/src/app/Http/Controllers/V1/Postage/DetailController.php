<?php

namespace App\Http\Controllers\V1\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Postage;
use CrsSdk\Shop\Store;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\Redis;

class DetailController extends BaseController
{
    /**
     * 邮费模板详情接口
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $data = (new Postage())->detail($input);
        if ($data['stores'] ?? []) {
            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => $data['stores'],
                'size'          => 10000,
            ])['list'];
            array_walk($storeList, static function ($value) use (&$list) {
                $list[] = Arr::only($value, ['id', 'name', 'short_name']);
            });
            $data['stores'] = $list ?? [];
        }
        return success($data);
    }
}
