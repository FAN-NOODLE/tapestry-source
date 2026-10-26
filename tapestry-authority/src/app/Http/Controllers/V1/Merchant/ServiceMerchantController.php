<?php

namespace App\Http\Controllers\V1\Merchant;

use App\Http\Controllers\V1\BaseController;
use App\Models\ServiceMerchantSetting;
use App\Services\EnterpriseService;
use App\Services\ServerMerchantService;
use Illuminate\Http\JsonResponse;

class ServiceMerchantController extends BaseController
{
    /**
     * @return JsonResponse
     */
    public function serviceMerchantDetail(): JsonResponse
    {
        $params = $this->request->all();

//        if (empty($params['enterprise_id']) && empty($params['hash']) && empty($params['id'])) {
//            error(700022);
//        }

        if (isset($params['merchant_id'])) {
            $merchant = (new ServerMerchantService())->getMerchantByMerchantId($params['merchant_id']);
            return success($merchant);
        }
        $serviceMerchantID = '';

        if (isset($params['enterprise_id'])) {
            $enterprise = (new EnterpriseService())->enterpriseDetail($params['enterprise_id']);
            $serviceMerchantID = $enterprise['service_merchant_id'];
        }

        if (isset($params['hash'])) {
            $enterprise = (new EnterpriseService())->detailByHash($params['hash']);
            $serviceMerchantID = $enterprise['service_merchant_id'];
        }

        if (isset($params['id'])) {
            $serviceMerchantID = $params['id'];
        }

        if (empty($serviceMerchantID)) {
            error(700023);
        }

        $merchant = (new ServerMerchantService())->merchantDetail($serviceMerchantID, ServiceMerchantSetting::ENABLE);

        return success($merchant);
    }


    public function serviceMerchantList()
    {
        return success((new ServerMerchantService())->merchantList());
    }

    public function getDefaultDetail()
    {
        return success((new ServerMerchantService())->defaultMerchant());
    }
}