<?php


namespace App\Http\Controllers\V1\WisdomGuide;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Common;

class BatchAddGuideController extends BaseController
{
    public function index()
    {
        $param = [
            'enterprise_id' => $this->enterpriseId,
            'qy_domain'     => $this->qyDomain
        ];
        Common::work()
            ->setApi('service/v1/store-worker/batch-add-guide')
            ->post($param);
        return success([]);
    }
}