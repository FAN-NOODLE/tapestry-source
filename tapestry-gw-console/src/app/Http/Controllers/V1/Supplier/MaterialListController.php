<?php


namespace App\Http\Controllers\V1\Supplier;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\MicroGoodsV1\Supplier;
use CrsSdk\Shop\Store;

class MaterialListController extends BaseController
{
    use ExcelTrait;

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '物料编码',
        '库存数量'
    ];

    /**
     * 门店的物料列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Supplier())->materialList($input);

        if ($input['source_id'] == $this->enterpriseId) {
            $storeName = '商城';
        } else {
            $storeName = (new Store())->detail([
                'id'            => $input['source_id'],
                'enterprise_id' => $input['enterprise_id']
            ])['name']; // short_name
        }

        $list['store_name'] = $storeName;

        return success($list);
    }

    /**
     * 导出门店的物料列表
     *
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function downloadMaterial()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['size'] = 999999;

        $list = (new Supplier())->materialList($input)['list'];

        $materialNoList = [];
        foreach ($list as $item) {
            $materialNoList[] = [
                'material_no' => $item['material_no'],
                'stock'       => $item['stock'],
            ];
        }

        $this->excelExport(self::EXCEL_TITLE, $materialNoList);
    }
}