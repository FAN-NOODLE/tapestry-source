<?php

namespace App\Http\Controllers\V1\Qywx\Channel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Channel;

class StatisticsListController extends BaseController
{
    /**
     * 渠道统计列表
     * {domain}/app/v1/qywx/channel/statistics-list
     *
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'start_time' => 'integer',
            'end_time'   => 'integer',
            'type'       => 'in:1,2,3',
            'page'       => 'integer',
            'size'       => 'integer'
        ], [
            'start_time.integer' => '[参数]start_time需为integer',
            'end_time.integer'   => '[参数]end_time需为integer',
            'type.in'            => '[参数]type值不合法',
            'page.integer'       => '[参数]page需为integer',
            'size.integer'       => '[参数]size需为integer'
        ]);
        $input = $this->request->only(['start_time', 'end_time', 'type', 'page', 'size', 'enterprise_id']);

        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Channel())->statisticsList($input);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }
}
