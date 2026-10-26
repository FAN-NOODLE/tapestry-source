<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/8/20
 * Time: 10:37
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\Work\BaseClient;
use Illuminate\Http\Request;
use CrsSdk\Work\DataCompass;
use Maatwebsite\Excel\Facades\Excel;

class RankController extends BaseController
{
    const LEVEL_KEY = 'area_level_key';

    const RANK_LENGTH = 100;

    const RANK_TYPE_COLUMN_MAP = [
        1 => 'cumulative_customer_num',
        2 => 'today_customer_num',
        5 => 'chat_cnt',
        6 => 'reply_percentage',
        7 => 'cumulative_binding_customer_num',
        8 => 'today_binding_customer_num'
    ];

    /**
     * 顾客数据排行榜
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $this->setParams($request, self::RANK_LENGTH);
        if (empty($params['rank_column'])) {
            return success(['list' => []]);
        }

        $rankList = (new DataCompass())->customerStatisticRankForApp($params);
//        $rankList = (new DataCompass())->customerStatisticRankForConsole($params);

        //获取企业导购区域配置
        $enterpriseConfig = (new EnterpriseConfig())->info([
            'enterprise_id' => $this->enterpriseId,
            'key'           => self::LEVEL_KEY
        ])['data'];
        $areaLevel        = 0;
        if (!empty($enterpriseConfig)) {
            $areaLevel = current(array_column($enterpriseConfig, 'value'));
            if ($areaLevel != 0) {
                foreach ($rankList as &$val) {
                    $pathArr                = explode('/', $val['dp_name_path']);
                    $val['area_level_list'] = array_slice($pathArr, 0 - $areaLevel);
                }
            }
        }
        return success([
            'area_level_config' => $areaLevel,
            'list'              => $rankList['list']
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
            $result = (new DataCompass())->customerStatisticRankForConsole($params);
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
//            $areaName = $value['dp_name'] ?? '';
//            if (!empty($value['area_level_list']) && count($value['area_level_list']) > 1) {
//                $areaName = ($value['area_level_list'][2] ?? '') . '/' . ($value['area_level_list'][1] ?? '') . '/' . ($value['area_level_list'][0] ?? '');
//                $areaName = trim($areaName, '//');
//            }
            $cellData[] = [
                $value['name'],
                $value['num'],
                $value['dp_name'],
                //                $areaName,
                ''
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
     *
     * @return array|string|null
     */
    public function setParams(Request $request, $pageSize = 0)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->enterpriseAdmin['mobile'] ?? '';
        $params['qy_domain']     = $this->qyDomain;
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['rank_column']   = !empty($params['type']) && isset(self::RANK_TYPE_COLUMN_MAP[$params['type']]) ? self::RANK_TYPE_COLUMN_MAP[$params['type']] : '';
        $params['rank_length']   = $pageSize;

        return $params;
    }
}