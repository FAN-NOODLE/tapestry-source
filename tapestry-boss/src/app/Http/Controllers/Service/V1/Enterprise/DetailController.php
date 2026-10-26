<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Models\BrandAuth;
use App\Models\Enterprise;
use Illuminate\Http\Request;
use App\Http\Resources\EnterpriseResource;
use App\Http\Controllers\Service\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 企业详情
     * @param Request $request
     * @param Enterprise $enterpriseModel
     * @return array|void
     */
    public function index(Request $request, Enterprise $enterpriseModel)
    {
        $params = $request->only(['enterprise_id', 'hash', 'corp_id', 'domain']);

        // 查询信息不能为空
        if (isset($params['enterprise_id']) && !empty($params['enterprise_id'])) {
            $enterprise = $enterpriseModel->query()->where('id', (string)$params['enterprise_id'])->first();

        } elseif (isset($params['hash']) && !empty($params['hash'])) {
            $enterprise = $enterpriseModel->query()->where('hash', (string)$params['hash'])->first();

        } elseif (isset($params['corp_id']) && !empty($params['corp_id'])) {
            $enterprise = $enterpriseModel->query()->where('corp_id', (string)$params['corp_id'])->first();

        } elseif (isset($params['domain']) && !empty($params['domain'])) {
            $enterprise = $enterpriseModel->query()->where('domain', (string)$params['domain'])->first();

        } else {
            return error(201102);
        }

        // 企业信息查询
        if (!$enterprise) {
            return error(201002);
        }

        $enterpriseArr = $enterprise->toArray();

        // 获取定制化权限标识
        $auth = (new BrandAuth())->getBrandAuth(
            $enterpriseArr['id'],
            '',
            [],
            [
                'enterprise_id',
                'brand_code'
            ]
        );
        $auth = array_column($auth, 'brand_code', 'enterprise_id');
        $enterpriseArr['brand_code'] = $auth && isset($auth[$enterpriseArr['id']]) ? $auth[$enterpriseArr['id']] : '';

        $enterpriseArr['identifier'] = 'contact';
        if ($enterpriseArr['contact_status'] == 3) {
            $enterpriseArr['identifier'] = 'crs-shop';
        }
        return success(['data' => $enterpriseArr]);
    }
}