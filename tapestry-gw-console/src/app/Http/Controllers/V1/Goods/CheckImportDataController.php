<?php
/**
 * 核对批量修改商品的数据
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use App\Micros\GoodsMicro;
use App\Traits\ImportTrait;
use Box\Spout\Reader\Common\Creator\ReaderEntityFactory;
use Illuminate\Http\Request;
use App\Common\Constant;

class CheckImportDataController extends BaseController
{
    use ImportTrait;

    public function index(Request $request)
    {
        $key = $this->getKey('update_goods');

        if (!$request->hasFile('file')) {
            notice(451506);
        };

        $params   = $this->getParams($request);
        $sourceId = (int)$params['source_id'];

        // 临时绝对路径
        $filePath = $request->file('file')->getRealPath();

        $data = $this->getData($filePath);

        $header = array_shift($data);

        if (empty($header) || (array_values($header) != Constant::BATCH_UPDATE_GOODS_HEADER)) {
            notice(451500);
        }

        if (count($data) > Constant::BATCH_UPDATE_LIMIT) {
            notice(451507);
        }

        $skuIds = array_column($data, 'sku_id');

        $goodsParams = [
            'enterprise_id' => $this->enterpriseId,
            'page'          => Constant::PAGE,
            'size'          => Constant::BATCH_UPDATE_LIMIT,
            'type'          => 1,
            'search'        => [
                'sku_id_list' => $skuIds,
                'source_ids'  => [$sourceId]
            ]
        ];

        $channelData     = (new GoodsMicro($goodsParams))->getChannelList();
        $channelList     = $channelData['list'] ?? [];
        $channelListById = array_column($channelList, null, 'sku_id');
        $channelSkuIds   = array_column($channelList, 'sku_id');

        $existedSkuIds = [];
        foreach ($data as $v) {
            array_walk($v, function (&$val, $key) {
                $val = trim($val);
            });
            $reasons = [];

            if (empty($v['sku_id'])) {
                $reasons[] = '商品不存在';
            }

            $skuId    = $v['sku_id'] ?? '';
            $skuName  = $v['sku_name'] ?? '';
            $skuPrice = $v['sku_price'] ?? '';

            if (in_array($skuId, $existedSkuIds)) {
                $reasons[] = '门店商品ID重复操作';
            }

            if (!in_array($skuId, $channelSkuIds)) {
                $reasons[] = '门店商品ID不存在';
            }

//            if (empty($skuName)) {
//                $reasons[] = '门店商品名称为空';
//            }

            if (!is_numeric($skuPrice) || $skuPrice <= 0) {
                $reasons[] = '门店商品价格有误';
            }

            if (empty($channelListById[$skuId])) {
                $reasons[] = '渠道商品不存在';
            }

            $channelInfo = $channelListById[$skuId] ?? [];

            $status = $channelInfo['status'] ?? 0;
            if ($status == 2) {
                $reasons[] = '门店商品上架状态不能编辑';
            }

            if (!empty($reasons)) {
                $v['reasons'] = $reasons;
                $this->pushData($v, $key, $this->errorType);
                continue;
            }

            $existedSkuIds[] = $skuId;
            $this->pushData($v, $key, $this->successType);
        }

        return success([
            'key'           => $key,
            'success_total' => $this->getDataTotal($key, $this->successType),
            'error_total'   => $this->getDataTotal($key, $this->errorType),
        ]);
    }

    protected function getData($filePath)
    {
        $data = [];
        try {
            $reader = ReaderEntityFactory::createXLSXReader();
            $reader->open($filePath);
            $reader->setShouldFormatDates(true);
            foreach ($reader->getSheetIterator() as $sheet) {
                foreach ($sheet->getRowIterator() as $i => $row) {

                    $cells   = $row->getCells();
                    $rowData = array_fill_keys(['sku_id', 'sku_name', 'sku_price'], '');

                    foreach ($cells as $k => $cell) {
                        $cellValue = $cell->getValue();

                        switch ($k) {
                            case 0:
                                $rowData['sku_id'] = $cellValue;
                                break;
                            case 1:
                                $rowData['sku_name'] = $cellValue;
                                break;
                            case 2:
                                $rowData['sku_price'] = (float)$cellValue * 100;
                                break;
                        }
                    }

                    $data[] = $rowData;
                }

                break;
            }
            $reader->close();
        } catch (\Exception $e) {
            notice(451506, [], [$e->getMessage()]);
        }

        return $data;
    }

    private function getParams(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'source_id' => 'required|integer',
        ];

        // 提示语
        $messages = [
            'source_id.*' => '渠道参数必填',
        ];

        $this->validate($request, $rules, $messages);

        return $input;
    }

}
