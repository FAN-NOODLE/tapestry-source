<?php


namespace App\Http\Controllers\V1\NewTemplate;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\NewTemplate;
use Exception;
use Illuminate\Http\JsonResponse;

class OnlineTplVersionController extends BaseController
{
    /**
     * 获取门店线上使用的模板
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $result = (new NewTemplate())->OnlineTplVersion([
            'enterprise_id' => $this->enterpriseId
        ]);

        $data = [];
        if (!empty($result)) {
            foreach ($result as $key => $val) {
                $data[] = [
                    'id'            => $val['id'],
                    'template_type' => $val['template_type'],
                ];
            }
        }

        return success($data);
    }
}
