<?php
/**
 * 根据字段查询企业ID集合
 * User: 何岳龙
 * Date: 2019/02/20
 * Time: 15:26
 */

namespace App\Http\Controllers\Service\V1\EnterpriseConfig;

use App\Models\EnterpriseConfig;
use App\Http\Controllers\Service\V1\BaseController;

class SelectEnterpriseIdController extends BaseController
{

    public function index()
    {
        $params = request()->only(['key', 'value']);

        if (empty($params['key'])) {
            return notice(620000);
        }

        $data = EnterpriseConfig::query()
            ->where('key', $params['key'])
            ->when(isset($params['value']), function ($query) use ($params) {
                return $query->where('value', $params['value']);
            })
            ->get(['enterprise_id'])
            ->toArray();

        $enterpriseIds = [];
        if (!empty($data)) {
            $enterpriseIds = array_column($data, 'enterprise_id');
        }

        return success(['data' => ['enterprise_ids' => $enterpriseIds]]);
    }
}
