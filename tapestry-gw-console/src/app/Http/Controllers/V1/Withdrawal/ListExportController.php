<?php


namespace App\Http\Controllers\V1\Withdrawal;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\WithdrawalListExportJob;
use Illuminate\Http\Request;

class ListExportController extends BaseController
{
    /*
     * 文件名
     */
    const FILE_NAME = '提现审批';

    const  STATUS = [
        1 => '待审批',
        2 => '提现成功',
        3 => '提现失败',
        4 => '已驳回',
        5 => '转账中'
    ];
    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'withdrawal/export_list';

    /**
     * 提现审批数据tag
     */
    const TAG = 'export-withdrawal-list';

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

        $name = self::FILE_NAME.'_'.self::STATUS[$input['status']].'_'.$time;

        // 导出
        (new WithdrawalListExportJob($input, $input, $name))->handle();

        return success();

    }

}