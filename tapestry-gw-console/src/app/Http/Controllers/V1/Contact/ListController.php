<?php

namespace App\Http\Controllers\V1\Contact;

use CrsSdk\Work\Contact;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 联系我二维码列表
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['page', 'size', 'channel', 'select_type']);

        $params['qy_domain']      = $this->qyDomain;
        $params['app_identifier'] = $this->appIdentifier;
        $params['enterprise_id']  = $this->enterpriseId;

        $list = (new Contact())->contactList($params);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }
}
