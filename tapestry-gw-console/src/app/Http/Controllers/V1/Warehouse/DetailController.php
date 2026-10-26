<?php
/**
 * Created by PhpStorm.
 * User: pw
 * Date: 2020/4/18
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Warehouse;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Warehouse;
use Illuminate\Support\Facades\Log;

class DetailController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input                  = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $data = (new Warehouse())->detail($input);
        } catch (Error $error) {

            Log::info('获取仓库详情失败', [$error->getMessage()]);
            notice(1621002, [], [$error->getMessage()]);

        } catch (\Exception $exception) {

            Log::info('获取仓库详情失败' . $exception->getMessage(), $input);
            notice(1621002);
        }

        // 返回结果
        return success($data);
    }

}