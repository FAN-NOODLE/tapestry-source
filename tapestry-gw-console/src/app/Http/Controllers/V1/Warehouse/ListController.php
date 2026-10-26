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

class ListController extends BaseController
{

    public function index()
    {
        $input                  = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $data = (new Warehouse())->lists($input);
        } catch (Error $error) {

            Log::info('获取仓库列表失败', $error->data);
            notice(1621001, [], $error->data);

        } catch (\Exception $exception) {

            Log::info('获取仓库列表失败' . $exception->getMessage(), $input);
            notice(1621001);
        }

        //返回结果
        return success($data);
    }
}