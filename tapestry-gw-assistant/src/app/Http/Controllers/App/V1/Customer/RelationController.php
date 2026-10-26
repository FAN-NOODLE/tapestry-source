<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\JsonResponse;

class RelationController extends BaseController
{
    /**
     * @var Sdk
     */
    protected $microScrm;

    public function __construct()
    {

        parent::__construct();
        $this->microScrm = Common::microScrm();
    }

    public function index():JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/get-user-member-relation')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}