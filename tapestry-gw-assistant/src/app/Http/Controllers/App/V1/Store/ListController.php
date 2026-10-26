<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/9
 * Time: 17:04
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\MicroScrm\GoMicroScrm;

use Illuminate\Http\JsonResponse;

/**
 * 门店列表
 * Class ListController
 * @package App\Http\Controllers\App\V1\Store
 */
class ListController extends Controller
{
    /**
     * 门店列表
     * @return JsonResponse
     */

    public function index(): JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new GoMicroScrm())->storeList($params);

        return success(['list' => $result]);
    }
}
