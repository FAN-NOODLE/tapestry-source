<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\AppManage;

use App\Services\EnterpriseAppService;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class ListController extends BaseController
{

    /**
     * 应用列表
     * @author 何岳龙 heyuelong@vchangyi.com
     * @date 2018/12/25 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->get('enterprise_id', 0);
/*        if (empty($enterpriseId)) {
            return error(201101);
        }*/

        // 查询企业应用列表
        $enterpriseService = new EnterpriseAppService();
        $data              = $enterpriseService->enterpriseAppList($enterpriseId);

        return success([
            'data' => $data
        ]);
    }
}
