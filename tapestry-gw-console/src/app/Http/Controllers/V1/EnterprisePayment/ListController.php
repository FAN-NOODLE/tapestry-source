<?php
/**
 * 获取企业微信红包活动列表
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    /**
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接值
        $params = $request->only(['search', 'page', 'size']);
        $params['search']['enterprise_id'] = $this->enterpriseId;
        $params['search']['qy_domain'] = $this->qyDomain;
        $params['search']['app_identifier'] = $this->appIdentifier;

        $data= (new EnterprisePayment())->lists($params);

        return success([
            'list'  => $data['data'],
            'page'  => $data['meta']['page'],
            'size'  => $data['meta']['size'],
            'total' => $data['meta']['total'],
        ]);
    }
}
