<?php


namespace App\Http\Controllers\Service\V1\FastRegisterWeApp;

use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Http\Request;
use App\Models\RegisterProgramAgent;

class DetailController extends BaseController
{


    public function index(Request $request, RegisterProgramAgent $model)
    {

        $params = $request->only([
            'id',
            'enterprise_id',
        ]);

        if (empty($params['enterprise_id'])) return error(600001);
        $info = $model->getInfoByEnterpriseId($params['enterprise_id']);
        $res = !empty($info) ? [$info->toArray()] : [];
        return success($res);
    }

}
