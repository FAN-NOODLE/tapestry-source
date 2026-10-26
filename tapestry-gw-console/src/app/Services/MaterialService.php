<?php


namespace App\Services;


use App\Common\Constant;
use CrsSdk\MicroGoodsV1\MaterialChannel;

class MaterialService extends BaseService
{

    /**
     * 渠道上下架状态 2-上架 4-暂停
     */
    const UP_STATUS = 2;
    const STOP_STAtUS = 4;

    /**
     * 获取物料id对应的库存数据数组
     * @param $enterpriseId
     * @param $materialId
     * @return array
     * @throws \Exception
     */
    public function getTotalStockArray($enterpriseId, $materialId)
    {
        $channelMaterialList = (new MaterialChannel)->MaterialChannelList([
            'enterprise_id' => $enterpriseId,
            'material_id'   => [$materialId],
            'statuss'       => [self::UP_STATUS, self::STOP_STAtUS],
            'page'          => Constant::PAGE,
            'size'          => Constant::MAX_DEFAULT_SIZE
        ]);

        $channellist = (new MaterialChannel())->sourceStockList([
            'enterprise_id' => $enterpriseId,
            'material_id'   => $materialId,
            'page'          => Constant::PAGE,
            'size'          => Constant::MAX_DEFAULT_SIZE
        ])['list'];

        // 2.处理同material_id下的库存合并
        $stockArr = [];
        if (isset($channelMaterialList['list']) && !empty($channelMaterialList['list'])) {
            foreach ($channelMaterialList['list'] as $key => $value) {
                $stockArr[$value['material_id']][] = $value['on_line_current_stock'] ?? 0;
                $stockArr[$value['material_id']][] = $value['off_line_current_stock'] ?? 0;
            }
        }

        return $data = [
            'stock_arr'     => $stockArr,
            'total_channel' => count($channellist),
        ];
    }

}