<?php


namespace App\Http\Controllers\App\V1\Customer;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\GoMicroScrm\Customer;
use Illuminate\Http\JsonResponse;

class LostRecordController extends BaseController
{

    public function index(): JsonResponse
    {
        $input                  = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = strval($this->mobile);
        try {
            $result = (new Customer())->lostRecord($input);
            return success($result['data'] ?? $result);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
    }

    public function option(): JsonResponse
    {
        $input                  = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = strval($this->mobile);
        try {
            $result = (new Customer())->lostOption($input);
            return success($result['data'] ?? $result);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
    }
}