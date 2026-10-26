<?php


namespace App\Http\Controllers\V1\DataCompass;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\DataExportByWorkerJob;
use CrsSdk\Work\DataCompass;
use Illuminate\Http\Request;

class DataExportByWorkerController extends BaseController
{

    /*
     * 文件名
     */
    const FILE_NAME = '导购顾客数据统计按导购';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'data_compass/export_by_worker';

    /**
     * 导购顾客数据统计tag
     */
    const TAG = 'export-by-worker';

    public function index(Request $request)
    {

        $input = $this->request->input();
        $input['for_page'] = false;
        $input['size']          = Constant::MAX_DEFAULT_SIZE;
        $input['page']          = Constant::PAGE;
        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain']     = $this->qyDomain;
        $input['operate_name']  = $this->enterpriseAdmin['username'] ?? '';
        $input['operate_id']    = $this->enterpriseAdminId;
        $input['export_type']   = $input['export_type'] ?? 1;
        $input['tag']           = $input['tag'] ?? self::TAG;
        $input['route_url']     = self::ROUTE_URL;
        $input['version']       = $input['version'] ?? '';
        $time                   = date('YmdHis', time());

        $name = self::FILE_NAME.$time;

        // 导出
        (new DataExportByWorkerJob($input, $input, $name))->handle();

        return success();

    }
}