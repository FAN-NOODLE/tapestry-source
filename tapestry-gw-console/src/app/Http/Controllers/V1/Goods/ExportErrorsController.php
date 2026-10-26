<?php
/**
 * 错误数据下载
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Traits\ImportTrait;
use Illuminate\Http\Request;
use Maatwebsite\Excel\Facades\Excel;
use function Complex\cot;

class ExportErrorsController extends BaseController
{
    use ImportTrait;
    protected $colWidth = [
        'A' => 25, 'B' => 25, 'C' => 25, 'D' => 30
    ];

    public function index(Request $request)
    {
        $ext      = Constant::BATCH_UPDATE_GOODS_EXT;
        $fileName = Constant::IMPORTED_ERROR_GOODS_FILENAME;
        $header   = Constant::IMPORTED_ERROR_GOODS_HEADER;
        $params   = $this->getParams($request);
        $key      = $params['key'];

        $list = $this->getDataList($key, $this->errorType);

        $data = [];
        foreach ($list as $v) {
            if (!isset($v['sku_id'])) {
                continue;
            }
            $data[] = [
                $v['sku_id'] . "\t",
                (string)$v['sku_name'],
                (string)($v['sku_price'] / 100),
                implode(PHP_EOL, $v['reasons']),
            ];
        }

        return Excel::create(date('YmdHis') . mt_rand(100, 999), function ($excel) use ($header, $data, $fileName) {
            $excel->sheet($fileName, function ($sheet) use ($header, $data) {
                $sheet->setWidth($this->colWidth)
                    ->rows([$header])
                    ->rows($data)
                    // 设置表头样式
                    ->cells('A1:D1', function ($cells) {
                        $cells
                            ->setAlignment('center')
                            ->setValignment('center')
                            ->setFontSize(13);
                    })
                    // 设置内容样式
                    ->cells("A2:D" . (count($data) + 1), function ($cells) {
                        $cells->setAlignment('left');
                    });
            });

        })->export($ext);

    }

    private function getParams(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'key' => 'required',
        ];

        // 提示语
        $messages = [
            'key.*' => '键值必填',
        ];

        $this->validate($request, $rules, $messages);

        return $input;
    }
}
