<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ShareChain\StaffShareStrategy;
use App\Services\ShareChain\Substance;
use App\Services\StoreWorkerService;
use Exception;

class DeleteGuideController extends BaseController
{

    public function index()
    {
        $rules = [
            'id'            => 'required',
            'operator_id'   => 'integer',
            'enterprise_id' => 'integer',
        ];
        $message = [
            'id.required' => 'id参数必填',
        ];
        $input = $this->validate($this->request, $rules, $message);
        try {
            StoreWorkerService::deleteGuide($input);
            // 更新用户分享记录
            $staffShareStrategy = new Substance(new StaffShareStrategy());
            $staffShareStrategy->unbindingLeave($input['enterprise_id'], [], [$input['id']], Constant::GUIDE_TYPE);
        } catch (Exception $e) {

        }
        return success();
    }

}
