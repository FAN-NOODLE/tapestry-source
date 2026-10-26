<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/10 0010
 * Time: 14:44
 */

namespace App\Http\Controllers\App\V1\Common;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Open\Logistics;

class LogisticsListController extends Controller
{
    protected $isRequireLogin = false;

    /**
     * 物流列表接口
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        try {
            $data = (new Logistics())->keyList()['data'];
        } catch (\Throwable $e) {
            error('651005');
        }

        $new = [];
        foreach ($data as $key => $value) {
            $new[] = ['key' => $key, 'value' => $value];
        }

        return success($new);
    }

}
