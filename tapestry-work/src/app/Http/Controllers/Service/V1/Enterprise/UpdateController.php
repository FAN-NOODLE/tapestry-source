<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use Illuminate\Http\Request;

class UpdateController extends BaseController
{
    /**
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $corpId = $request->post('corp_id');
        $name = $request->post('name');
        $qyDomain = $request->post('qy_domain');

        if (empty($corpId)) {
            throw new Error(100000, 'corp_id 字段不能为空');
        }

        $params = [
            'corpId' => $corpId,
            'epName' => $name,
            'epDomain' => config('client')['api_url']['work'],
            'corpName' => $name
        ];

        $result = (new Enterprise())->enterpriseUpdate($qyDomain, $params);
        return $result;
    }
}
