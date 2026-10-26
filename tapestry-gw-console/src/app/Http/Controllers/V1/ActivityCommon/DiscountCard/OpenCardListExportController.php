<?php


namespace App\Http\Controllers\V1\ActivityCommon\DiscountCard;


use App\Http\Controllers\V1\BaseController;
use App\Jobs\Activity\OpenCardListExportJob;

class OpenCardListExportController extends BaseController
{
    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = '先享卡-开卡信息';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'activity-common/discount-card/open-card-export';

    /**
     * 列表出tag
     */
    const TAG = 'open-card-list';

    public function index()
    {
        $input = $this->request->only([
            'rule_id',
            'search',
            'tag',
            'route_url'
        ]);

        $params                  = $input;
        $params['tag']           = $input['tag'] ?? self::TAG;
        $params['route_url']     = self::ROUTE_URL;
        $params['operate_name']  = $this->enterpriseAdmin['username'];
        $params['operate_id']    = $this->enterpriseAdminId;
        $params['enterprise_id'] = $input['enterprise_id'] = $this->enterpriseId;
        $params['export']        = $input['export'] = 1; // 1-导出 2-列表

        // job开始执行
        $job = new OpenCardListExportJob($params, $input, self::FILE_NAME);

        //$job->handle();
        $this->dispatch($job);

        return success([]);
    }

}