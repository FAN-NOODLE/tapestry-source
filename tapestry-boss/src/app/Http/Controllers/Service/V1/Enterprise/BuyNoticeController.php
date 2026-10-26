<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseBuyRecord;

class BuyNoticeController extends BaseController
{
    /**
     * 购买提醒
     * @return array
     */
    public function index()
    {
        $enterpriseId = $this->request->input('enterprise_id');

        $buyRecord = EnterpriseBuyRecord::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('end_time', '>', time())
            ->where('end_time', '<=', time() + 30 * 24 * 60 * 60)
            ->orderBy('end_time', 'asc')
            ->first();

        if (empty($buyRecord)) {
            return [
                'user_count' => 0,
                'days'       => 0,
            ];
        }

        return [
            'user_count' => $buyRecord['user_count'],
            'days'       => (int)(1 + ($buyRecord['end_time'] - time()) / (24 * 60 * 60))
        ];
    }
}