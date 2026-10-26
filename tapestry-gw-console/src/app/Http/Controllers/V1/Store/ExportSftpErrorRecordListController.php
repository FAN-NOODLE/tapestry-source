<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\Sdk;

class ExportSftpErrorRecordListController extends BaseController
{
    use ExcelTrait;

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $res = Sdk::work()->setApi('service/v1/store/sftp-error-record-list')
            ->post($params);
        if (empty($res['error_record'])) {
            return success([]);
        }
        $dataArr = [];
        foreach ($res['error_record'] as $str) {
            $arr = json_decode($str, true);
            if (!empty($arr)) {
                $dataArr[] = $arr;
            }
        }
        if (count($dataArr[0]) >= 10) {
            $header = ['MobileNumber','Name','AliasName','EmployeeId',
                'JobLocation','JobPosition','Gender','EmailAddress',
                'Status','StoreCode','fail_reason','is_success'];

        } else {
            $header = ["Department_ID",'fail_reason','is_success'];
        }
        $this->excelExport($header, $dataArr, $res['file_name']);
    }

}
