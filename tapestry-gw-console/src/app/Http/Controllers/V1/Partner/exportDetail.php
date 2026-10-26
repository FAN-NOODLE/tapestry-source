<?php


namespace App\Http\Controllers\V1\Partner;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\PartnerDetailExport;
use Illuminate\Http\Request;

class exportDetail extends BaseController
{


    /*
     * 文件名
     */
    const FILE_NAME = '合伙人余额明细';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'partner/export_detail';

    /**
     * 合伙人列表tag
     */
    const TAG = 'export-partner-detail';

    public function index(Request $request)
    {

        $input = $this->request->input();
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
        (new PartnerDetailExport($input, $input, $name))->handle();

        return success();

    }

}