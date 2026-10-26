<?php

namespace App\Http\Controllers\V1\Contact;

use CrsSdk\Work\Contact;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SelectionListController extends BaseController
{
    /**
     * 联系我二维码列表(下拉框选择)
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params['qy_domain']       = $this->qyDomain;
        $params['app_identifier']  = $this->appIdentifier;
        $params['enterprise_id']   = $this->enterpriseId;
        $params['size']            = 0; // 不分页
        $params['has_welcome_msg'] = 2; // 是否包含已经添加欢迎语的联系我二维码：1、是，2、否

        $list = (new Contact())->contactList($params);

        return success(['list' => $list['data']]);
    }
}
