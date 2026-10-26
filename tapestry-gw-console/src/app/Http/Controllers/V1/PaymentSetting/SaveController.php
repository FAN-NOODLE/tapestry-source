<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\PaymentSetting;

use CrsSdk\Shop\Payment;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $params = $request->only(['merchant_id', 'merchant_key', 'api_v3_key', 'serial_no']);

        $privateFile = $request->file('private_key');
        $publicFile  = $request->file('public_key');

        $params['enterprise_id'] = $this->enterpriseId;

        if ($request->hasFile('private_key')) {
            if ($privateFile->getClientOriginalExtension() != 'pem') {
                error(410003);
            }
            $params['private_key'] = file_get_contents($request->private_key);
        }

        if ($request->hasFile('public_key')) {
            if ($publicFile->getClientOriginalExtension() != 'pem') {
                error(410003);
            }
            $params['public_key'] = file_get_contents($request->public_key);
        }

        $this->writeOperatorLog('编辑支付:' . '商户ID:' . $params['merchant_id']);

        $data = (new Payment())->save($params);

        return success($data);
    }

}
