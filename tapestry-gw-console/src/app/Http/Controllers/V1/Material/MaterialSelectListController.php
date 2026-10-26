<?php

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;

class MaterialSelectListController extends BaseController
{

    /**
     * 物料启用
     */
    const STATUS_START = 1;

    /**
     * 物料禁用
     */
    const STATUS_STOP  = 2;

    /**
     * 渠道的物料列表
     * 
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        $input = $this->request->all();

        try {
            $input['enterprise_id'] = $this->enterpriseId;
            $input['join_stock'] = 'all'; // 这个参数控制展示的物料是只有带供货关系的，all是所有
            $input['status'] = self::STATUS_START;
            // 超管source_id等于enterprise_id
            $input['source_id'] = $this->storeIds;
            if ($this->isSuperAdmin == 1) {
                $input['source_id'] = [$this->enterpriseId];
            }

            $data = (new Material())->assignChannelList($input);
            foreach ($data['list'] as $key=>$item) {
                $data['list'][$key] = collect($item)->only(['material_no', 'name', 'material_id', 'cost_price', 'sign_price'])->toArray();
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