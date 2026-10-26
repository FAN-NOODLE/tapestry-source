<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Jobs\StoreWorkerGenerateQrCodeJob;
use App\Models\StoreWorker;
use App\Services\MemberService;
use App\Services\StoreWorkerService;
use Exception;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class GenerateWorkerPosterController extends BaseController
{
    /**
     * 手动生成导购海报
     * @return array
     * @throws ValidationException
     * @throws Exception
     */
    public function index()
    {
        $rules   = [
            'enterprise_id' => 'required',
            'id'            => 'required',
            'domain'        => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业id参数必填',
            'id.required'            => '导购id参数必填',
            'domain.required'        => '企业domain参数必填',
        ];
        $input   = $this->validate($this->request, $rules, $message);
        $data = [
            'worker_ids'    => [$input['id']],
            'enterprise_id' => $input['enterprise_id'],
            'store_id'      => 0,
            'domain'        => $input['domain'],
            'operator_id'   => $input['operator_id'] ?? 0,
            'operator_name' => $input['operator_name'] ?? '系统'
        ];
        //同步生成二维码和太阳码
        //(new StoreWorkerGenerateQrCodeJob($data))->handle();
        dispatch(new StoreWorkerGenerateQrCodeJob($data));

        return success();
    }
}
