<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:48
 */

namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\FiscalMonth;
use Illuminate\Http\Request;

class FiscalMonthController extends BaseController
{

    /**
     *  当前时间标识
     */
    CONST FISCAL_MONTH_NOW = 1;
    /**
     * 财月获取
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request): array
    {
        $enterpriseId = $request->post('enterprise_id', 0);
        $list = FiscalMonth::query()->where('enterprise_id', $enterpriseId)
            ->get(['period', 'start_time', 'end_time', 'year', 'fiscal_year'])->toArray();
        $now = time();
        foreach ($list as &$val) {
            $val['now'] = 0;
            if ($val['start_time'] < $now && $val['end_time'] > $now) {
                $val['now'] = self::FISCAL_MONTH_NOW;
            }
        }
        return success($list);
    }
}