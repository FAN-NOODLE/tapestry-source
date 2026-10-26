<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/25
 * Time: 18:52
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;

class ShopGoodsListController extends BaseController
{
    public function index()
    {
        //查询条件校验
        $input = $this->request->all();

        if (empty($input['source_id'])) notice(100000, ['门店ID不能为空']);

        $input['page'] = $input['page'] ?? 1;
        $input['size'] = $input['size'] ?? 10;
        try {
            // 获取门店的物料操作权限
            $materialRule = (new Material())->assignChannelList([
                'enterprise_id' => $this->enterpriseId,
                'source_id'     => [$input['source_id']],
                'size'          => 10000,
                'status'        => 1
            ]);
            if (empty($materialRule['list'])) {
                return success($materialRule);
            }

            // 门店可以操作的物料
            $materialIds = array_column($materialRule['list'], 'material_id');

            $input['enterprise_id'] = $this->enterpriseId;
            $input['need_category'] = 1;
            $input['search']['source_ids'] = [$this->enterpriseId];
            $input['search']['material_ids'] = $materialIds;
            $input['search']['name_type'] = $input['name_type'] ?? '';
            $input['search']['keywords'] = $input['keywords'] ?? '';
            $input['type'] = 6; // 已发布的普通商品

            $skuInfo = (new Goods())->list($input);

            $skuInfo['list'] = $this->formatInfo($skuInfo['list']);;
        } catch (Error $error) {

            Log::info('获取商品列表失败', [$error->getMessage()]);
            notice(100000, [$error->getMessage()], $error->data);

        } catch (\Exception $exception) {

            Log::info('获取商品列表失败' . $exception->getMessage(), $input);
            notice(100000, [$exception->getMessage()]);
        }

        //返回结果
        return success($skuInfo);
    }


    protected function formatInfo($skuInfo)
    {
        $data = [];

        if (empty($skuInfo)) return $data;

        foreach ($skuInfo as $item) {
            $data[] = [
                'spu_id'        => $item['spu_id'],
                'sku_id'        => $item['sku_id'],
                'sku_name'      => $item['sku_name'],
                'sku_img'       => $item['sku_img'],
                'material_id'   => $item['material_id'],
                'category_list' => $item['category_list'] ?? [],
                'market_price'  => $item['info'][0]['price'] ?? 0,
            ];
        }

        return $data;

    }
}