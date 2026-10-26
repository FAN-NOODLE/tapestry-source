<?php

namespace App\Http\Controllers\App\V1\EnterpriseNew;

use App\Models\Industry;
use App\Models\Enterprise;
use App\Http\Controllers\App\V1\BaseController;

class ChangeBaseInfoController extends BaseController
{
    public function index()
    {
        $params = $this->request->input();

        if (empty($params['enterprise_id'])) {
            error(201108);
        }

        $enterprise = Enterprise::searchEnterprise($params['enterprise_id']);
        if (empty($enterprise)) {
            error(201108);
        }

        if (isset($params['industry_id']) && $params['industry_id'] !== '') {
            $exists = Industry::idExists($params['industry_id']);
            if ($exists === false) {
                error(720001);
            }

            $enterprise->update([
                'industry_id' => $params['industry_id'],
            ]);

            return success();
        }

        if (!empty($params['business_source'])) {
            // 1:企业微信,10:畅移（直签）,11:畅移（代理）,12:畅移（试用）20:畅鲟（直签）21:畅鲟（代理）22:畅鲟（试用）',
            if (!in_array($params['business_source'], [1, 10, 11, 12, 20, 21, 22])) {
                error(720002);
            }

            $enterprise->update([
                'business_source' => $params['business_source'],
            ]);
        }

        return success();
    }
}