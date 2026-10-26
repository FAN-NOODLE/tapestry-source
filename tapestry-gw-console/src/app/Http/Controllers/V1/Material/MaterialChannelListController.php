<?php

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class MaterialChannelListController extends BaseController
{
    /**
     * 渠道的物料列表
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->all();

        try {
            $input['enterprise_id'] = $this->enterpriseId;
            $input['join_stock'] = 'all'; // 这个参数控制展示的物料是只有带供货关系的，all是所有
            $data = (new Material())->assignChannelList($input);

            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => array_column($data['list'], 'source_id'),
                'size'          => 999999
            ])['list'];
            $storeList = array_column($storeList, 'name', 'id');
            foreach ($data['list'] as $key => $value) {
                $data['list'][$key]['store_name'] = $storeList[$value['source_id']] ?? '商城或其他';
            }

        } catch (Error $error) {

            Log::error('获取渠道的物料列表失败', $error->data);
            notice(450000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('获取渠道的物料列表失败' . $e->getMessage(), $input);
            notice(450000, $e->getMessage());

        }

        return success($data);
    }
}