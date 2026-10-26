<?php

namespace App\Http\Controllers\V1\Qywx\Channel;

use CrsSdk\Work\Channel;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class StatisticsDetailController extends BaseController
{
    /**
     * 渠道统计明细
     * @return \Illuminate\Http\JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'id'         => 'required|integer',
            'start_time' => 'integer',
            'end_time'   => 'integer',
            'type'       => 'in:1,2,3'
        ], [
            'id.required'        => '[参数]渠道id必填',
            'id.integer'         => '[参数]渠道id需为integer',
            'start_time.integer' => '[参数]start_time需为integer',
            'end_time.integer'   => '[参数]end_time需为integer',
            'type.in'            => '[参数]type值不合法'
        ]);
        $input = $this->request->only(['id', 'start_time', 'end_time', 'type']);

        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Channel())->statisticsDetail($input);

        return success($list['data']);
    }
}
