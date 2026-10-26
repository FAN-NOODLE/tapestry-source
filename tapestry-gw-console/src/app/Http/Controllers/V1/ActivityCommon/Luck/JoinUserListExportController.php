<?php


namespace App\Http\Controllers\V1\ActivityCommon\Luck;


use App\Http\Controllers\V1\BaseController;
use App\Jobs\Activity\LuckJoinUserExport;

class JoinUserListExportController extends BaseController
{
    // route: {domain}/gw-console/v1/activity-common/activity
    // type: luck
    // method: JoinUserListExport
    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = '好友助力抽奖参与用户';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'activity-common/activity';

    /**
     * 列表出tag
     */
    const TAG = 'luck-join-user-list';

    public function index()
    {
        $input = $this->request->only([
            'id',
            'tag',
            'route_url',
            'method'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['type']          = 1;

        $params                  = $input;
        $params['tag']           = $input['tag'] ?? self::TAG;
        $params['route_url']     = self::ROUTE_URL;
        $params['operate_name']  = $this->enterpriseAdmin['username'];
        $params['operate_id']    = $this->enterpriseAdminId;
        $params['enterprise_id'] = $input['enterprise_id'] = $this->enterpriseId;

        // job开始执行
        $job = new LuckJoinUserExport($params, $input, self::FILE_NAME.date("YmdHis", time()));
        // $job->handle();
         $this->dispatch($job);

        return success([]);
    }


}