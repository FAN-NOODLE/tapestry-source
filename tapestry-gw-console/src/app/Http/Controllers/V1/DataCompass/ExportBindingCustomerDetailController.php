<?php


namespace App\Http\Controllers\V1\DataCompass;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\ExportBindingCustomerDetailJob;
use CrsSdk\Work\StaffShare;
use Illuminate\Http\Request;

class ExportBindingCustomerDetailController extends BaseController
{

    const EXCEL_LIMIT = 100000;

    const PAGE_SIZE = 1000;

    /*
     * 文件名
     */
    const FILE_NAME = '导购绑客明细';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'data_compass/export_binding_customer_detail';

    /**
     * 导购顾客数据统计tag
     */
    const TAG = 'export-binding-customer-detail';

    public function index(Request $request)
    {

        $input = $this->request->input();
        $input['enterprise_hash'] = $this->request->header('enterprise-hash');
        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain']     = $this->qyDomain;
        $input['operate_name']  = $this->enterpriseAdmin['username'] ?? '';
        $input['operate_id']    = $this->enterpriseAdminId;
        $input['export_type']   = $input['export_type'] ?? 1;
        $input['tag']           = $input['tag'] ?? self::TAG;
        $input['route_url']     = self::ROUTE_URL;
        $input['version']       = $input['version'] ?? '';
        $input['is_page']       = true;
        $time                   = date('YmdHis', time());


        $queryTotalParams = $input;
        $queryTotalParams['size'] = 1;
        $queryTotalParams['page'] = 1;
        $totalResult = (new StaffShare())->bindingCustomerDetail($queryTotalParams);
        if ($totalResult['total'] > self::EXCEL_LIMIT) {
            $excelCount = (int) ceil($totalResult['total'] / self::EXCEL_LIMIT);
            for ($n = 1; $n <= $excelCount; $n++) {
                $input['size'] = self::PAGE_SIZE;
                $input['page'] = ($n - 1) * (self::EXCEL_LIMIT / self::PAGE_SIZE) + 1;
                if ($n == $excelCount) {
                    $input['page_total'] = ceil($totalResult['total']/ self::PAGE_SIZE);
                } else {
                    $input['page_total'] = (self::EXCEL_LIMIT / self::PAGE_SIZE) * $n;
                }
                $name = self::FILE_NAME . "($n)" . $time;
                dispatch(new ExportBindingCustomerDetailJob($input, $input, $name));
            }
        } else {
            $input['size'] = self::PAGE_SIZE;
            $input['page'] = 1;
            $input['page_total'] = ceil($totalResult['total'] / self::PAGE_SIZE);
            $name = self::FILE_NAME . $time;
            // 导出
            dispatch(new ExportBindingCustomerDetailJob($input, $input, $name));
        }

        return success();

    }

}