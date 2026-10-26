<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-31
 * Time: 15:55
 */

namespace App\Http\Controllers\V1\Member;

use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\CustomerExportJob;
use Illuminate\Http\Request;

class ExportMemberListController extends BaseController
{
    /**
     * 导出客户列表
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $all                   = $request->all();
        $all['qy_domain']      = $this->qyDomain;
        $all['app_identifier'] = $this->appIdentifier;
        $all['enterprise_id']  = $this->enterpriseId;
        $all['operate_id']     = $this->enterpriseAdminId;
        $all['operate_name']   = $this->enterpriseAdmin['username'];
        $all['tag']            = 'qy_member_list';
        $all['route_url']      = 'member/export-customer-list';

        $downloadName = "客户列表_" . date("YmdHi");
        try {
            dispatch(new CustomerExportJob($all, $all, $downloadName));
        } catch (\Exception $error) {
            notice($error->getCode(), $error->getMessage());
        }
        return success();
    }
}
