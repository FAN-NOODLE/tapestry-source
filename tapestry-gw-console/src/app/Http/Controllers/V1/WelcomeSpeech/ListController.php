<?php

namespace App\Http\Controllers\V1\WelcomeSpeech;

use Illuminate\Http\Request;
use CrsSdk\Work\WelcomeSpeech;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 欢迎语列表
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params                  = $request->only(['page', 'size', 'channel']);
        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new WelcomeSpeech())->list($params);

        return success([
            'list'    => $list['data'],
            'page'    => $list['meta']['page'],
            'size'    => $list['meta']['size'],
            'total'   => $list['meta']['total'],
            'default' => $list['default'],
        ]);
    }
}
