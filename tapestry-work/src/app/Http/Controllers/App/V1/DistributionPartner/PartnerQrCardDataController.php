<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/18
 * Time: 10:28 上午
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\AppCode;
use CrsSdk\Shop\Store;
use Illuminate\Http\Request;

class PartnerQrCardDataController extends BaseController
{
    protected $requireLogin = false;

    public function index(Request $request)
    {
        $storeId = $request->post('store_id');

        //TODO 换路径
        //太阳码参数
        $params = [
            'enterprise_id' => $this->enterpriseId ?? 2,
            'path'          => 'pages/load/load',
            'params'        => [
                'enterprise_id' => $this->enterpriseId ?? 119419154973376,
                'mem_uid'       => $this->memUid ?? 0,
                'store_id'      => $storeId
            ]
        ];
        try {
            $result = (new AppCode())->create($params);


        } catch (\Exception $exception) {
            throw new Error(600013, '生成海报失败' . $exception->getMessage());

        }
        return successful($result);

    }
}