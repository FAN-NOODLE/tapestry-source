<?php


namespace App\Http\Controllers\V1\WisdomGuide;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Micros\Common;
use Illuminate\Http\JsonResponse;

class UpgradeAdminController extends BaseController
{
    const TYPE_UPGRADE_ADMIN = 4;
    const TYPE_UPGRADE_GUIDE = 3;

    public function index(): JsonResponse
    {
        $result = [];
        if ($this->request->input('type') == self::TYPE_UPGRADE_ADMIN) {
            $this->validate($this->request, [
                'type'       => 'required',
                'name'       => 'required',
                'mobile'     => 'required',
                'range'      => 'required|array',
                'range_type' => 'required'
            ]);

            $param  = [
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->request->input('mobile'),
                'range'         => $this->reformRange(),
                'type'          => $this->request->input('type'),
                'range_type'    => $this->request->input('range_type'),
                'name'          => $this->request->input('name'),
                'operator_id'   => $this->enterpriseAdminId
            ];
            $result = Common::work()->setApi('service/v1/store-worker/upgrade-admin')->post($param);
        }
        if ($this->request->input('type') == self::TYPE_UPGRADE_GUIDE) {
            $param  = [
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->request->input('mobile'),
                'type'          => $this->request->input('type'),
                'range_type'    => $this->request->input('range_type'),
                'operator_id'   => $this->enterpriseAdminId
            ];
            $result = Common::work()->setApi('service/v1/store-worker/upgrade-admin')->post($param);
        }

        return success($result);

    }


    private function reformRange(): array
    {
        $range = $this->request->input('range');
        if (!$range) {
            throw new Error(10000, '请传入部门或门店id');
        }
        $reformRange = [];
        foreach ($range as $objectId) {
            $reformRange[] = [
                "type"      => $this->request->input('range_type'),
                "object_id" => $objectId
            ];
        }
        return $reformRange;
    }


}