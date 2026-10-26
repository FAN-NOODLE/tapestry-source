<?php


namespace App\Http\Controllers\V1\ActivityCommon\Help;


use App\Http\Controllers\V1\BaseController;
use App\Jobs\Activity\HelpInitListExportJob;

class InitUserListExportController extends BaseController
{
    // route: {domain}/gw-console/v1/activity-common/activity
    // type: luck
    // method: InitRankListExport
    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = '好友助力%s发起用户';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'activity-common/activity';

    /**
     * 列表出tag
     */
    const TAG = 'help-init-list';

    public function index()
    {
        $input = $this->request->only([
            'id',
            'activity_type',
            'tag',
            'type',
            'method'
        ]);

        $text = '';
        $activityType = $input['activity_type'] ?? 0;
        // 付邮
        if ($activityType == 6) {
            $text = '付邮';
        }

        if ($activityType == 7) {
            $text = '买赠';
        }

        $fileName = sprintf(self::FILE_NAME, $text);

        $params                  = $input;
        $params['tag']           = $input['tag'] ?? self::TAG;
        $params['route_url']     = self::ROUTE_URL;
        $params['operate_name']  = $this->enterpriseAdmin['username'];
        $params['operate_id']    = $this->enterpriseAdminId;
        $params['enterprise_id'] = $input['enterprise_id'] = $this->enterpriseId;

        // job开始执行
        $job = new HelpInitListExportJob($params, $input, $fileName . date("YmdHis", time()));
        //$job->handle();
        $this->dispatch($job);

        return success([]);
    }

}