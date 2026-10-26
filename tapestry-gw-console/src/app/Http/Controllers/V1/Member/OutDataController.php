<?php

namespace App\Http\Controllers\V1\Member;

use CrsSdk\Work\Customer;
use Maatwebsite\Excel\Facades\Excel;
use App\Http\Controllers\V1\BaseController;

class OutDataController extends BaseController
{
    /**
     * 客户统计导出
     * @return mixed
     * @throws \Exception
     */
    public function outData()
    {
        $params                   = $this->request->all();
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;

        $data = (new Customer)->outData($params);

        $name = "客户统计导出" . date('YmdHi', time());

        return Excel::create($name, function ($excel) use ($data) {
            $excel->sheet('data', function ($sheet) use ($data) {
                $sheet->rows($data['data']);
            });
        })->export('xlsx');
    }
}