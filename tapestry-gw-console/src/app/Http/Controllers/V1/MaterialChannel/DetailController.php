<?php


namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MaterialChannel;

class DetailController extends BaseController
{

    /**
     * 渠道详情
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'material_channel_id' => 'required|integer',
            'source_id'           => 'required|integer',
        ], [
            'material_channel_id.*' => '渠道ID不能为空',
            'source_id.*'           => '门店ID不能为空'
        ]);

        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $detail = (new MaterialChannel())->detail($input);

        $detail = $this->formatData($detail);

        return success($detail);
    }

    private function formatData($detail)
    {
        $return = [];
        if (empty($detail)) return $return;

        $return = [
            'id' => $detail['id'],
            'enterprise_id' => $detail['enterprise_id'],
            'material_id' => $detail['material_id'],
            'source_id' => $detail['source_id'],
            'sku_id' => $detail['sku_id'],
            'on_line_stock' => $detail['on_line_current_stock'],
            'off_line_stock' => $detail['off_line_current_stock'],
        ];

        return $return;
    }
}