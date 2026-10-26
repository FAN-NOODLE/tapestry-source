<?php

namespace App\Http\Controllers\V1\Qywx\Channel;

use CrsSdk\Work\Channel;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class ListController extends BaseController
{
    /**
     * 渠道列表
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'page' => 'integer',
            'size' => 'integer'
        ], [
            'page.integer' => '[参数]page需为integer',
            'size.integer' => '[参数]size需为integer'
        ]);
        $input = $this->request->only(['channel', 'page', 'size']);

        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Channel())->channelList($input);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }
}
