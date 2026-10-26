<?php
/**
 * 下载批量修改商品模板
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use Maatwebsite\Excel\Facades\Excel;

class ExportTemplateController extends BaseController
{
    protected $requireLogin = false;
    protected $colWidth = [
        'A' => 25, 'B' => 25, 'C' => 25
    ];

    public function index()
    {
        $fileName = md5(Constant::BATCH_UPDATE_GOODS_FILENAME);

        $path = storage_path('download') . '/' . $fileName . '.' . Constant::BATCH_UPDATE_GOODS_EXT;

        if (!file_exists($path)) {
            $header = Constant::BATCH_UPDATE_GOODS_HEADER;

            Excel::create($fileName, function ($excel) use ($header) {
                $excel->sheet(Constant::BATCH_UPDATE_GOODS_FILENAME, function ($sheet) use ($header) {
                    $sheet->setWidth($this->colWidth)
                        ->rows([$header])
                        // 设置表头样式
                        ->cells('A1:C1', function ($cells) {
                            $cells
                                ->setAlignment('center')
                                ->setValignment('center')
//                            ->setFontWeight(true)
//                            ->setFontColor('#aabbcc')
                                ->setFontSize(13);
                        });
                });
            })->store(Constant::BATCH_UPDATE_GOODS_EXT, storage_path('download'));
        }

        return response()->download($path);
    }
}
