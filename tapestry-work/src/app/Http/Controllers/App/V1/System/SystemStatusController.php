<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/2/21
 * Time: 16:48
 */

namespace App\Http\Controllers\App\V1\System;


use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use App\Exceptions\Error;
use CrsSdk\Boss\Enterprise;

class SystemStatusController extends BaseController
{
    /**
     * 该接口非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 获取crs小程序和商城小程序的开通状态
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseHash = $request->header('Enterprise-Hash');

        if (empty($enterpriseHash)) {
            throw new Error(110002, '企业HASH不能为空');
        }

        // 获取该企业的授权信息
        $enterpriseData = (new Enterprise())->detail(['hash' => $enterpriseHash]);
        $enterpriseData = $enterpriseData['data'];
        $enterpriseId = $enterpriseData['id'];
        $crsStatus = $enterpriseData['crs_status'];
        $storeStatus = $enterpriseData['store_status'];

        return success(
            [
                'data' => [
                    'enterprise_id' => $enterpriseId,
                    'crs_status' => $crsStatus,
                    'store_status' => $storeStatus,
                ]
            ]
        );
    }
}
