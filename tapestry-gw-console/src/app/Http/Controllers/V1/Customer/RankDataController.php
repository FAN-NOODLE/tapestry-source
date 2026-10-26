<?php

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use CrsSdk\Boss\EnterpriseConfig;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Maatwebsite\Excel\Facades\Excel;

class RankDataController extends BaseController
{
    /**
     * 顾客数据排行榜
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request): JsonResponse
    {
        $params = $this->setParams($request);

        $rankList = GoMicroScrm::rank($params);
        $list = [];
        foreach ($rankList['list'] as $val) {
            $list[] = [
                'shopping_guide_name'   => $val['name'],
                'num'   => $val['total'],
                'belonging_store'   => $val['dp_name'],
                'belonging_area' => ''
            ];

        }
        return success([
            'area_level_config' => 0,
            'list' => $list,
            'total' => $rankList['total'],
        ]);
    }

    /**
     * 设置导出列宽
     *
     * @var array
     */
    protected $colWidth = [
        'A' => 15,
        'B' => 15,
        'C' => 15,
        'D' => 15,
    ];

    public function export(Request $request)
    {
        $params = $this->setParams($request);
        try {
            $result = (new GoMicroScrm())->rank($params);;
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $cellData[] = [
            '导购',
            '企微好友数',
            '门店',
//            '所属区域',
        ];
        foreach ($result['list'] ?? [] as $value) {

            $cellData[] = [
                $value['name'],
                $value['total'],
                $value['dp_name'],
//                ''
            ];
        }

        $name = '导购排行榜导出';
        return Excel::create($name, function ($excel) use ($cellData) {
            $excel->sheet('ba-top', function ($sheet) use ($cellData) {
                $sheet->setWidth($this->colWidth)->rows($cellData);
            });
        })->export('xlsx');
    }

    /**
     * @param Request $request
     * @return array|string|null
     */
    public function setParams(Request $request)
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;

        return $params;
    }
}