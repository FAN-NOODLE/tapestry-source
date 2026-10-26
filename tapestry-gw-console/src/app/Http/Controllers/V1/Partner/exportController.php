<?php


namespace App\Http\Controllers\V1\Partner;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\PartnerListExport;
use Illuminate\Http\Request;

class exportController extends BaseController
{

    /*
     * 文件名
     */
    const FILE_NAME = '合伙人列表';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'partner/export_list';

    /**
     * 合伙人列表tag
     */
    const TAG = 'export-partner-list';

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
        (new PartnerListExport($input, $input, $name))->handle();

        return success();

    }
}