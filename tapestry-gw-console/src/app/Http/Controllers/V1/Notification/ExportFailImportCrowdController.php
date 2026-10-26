<?php

namespace App\Http\Controllers\V1\Notification;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;

class ExportFailImportCrowdController extends BaseController
{
    use ExcelTrait;

    /**
     * User: zhaoyang@vchangyi.com
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index()
    {
        $failStoreCodes = $this->request->input('fail_store_codes', []);
        if (!empty($failStoreCodes)) {
            $failData = [];
            foreach ($failStoreCodes as $storeCode) {
                $failData[] = [$storeCode, '未匹配到门店'];
            }
            $header = ['*门店code', '失败数据原因'];
            $this->excelExport($header, $failData, 'storeCodeError');
        }
    }

}
