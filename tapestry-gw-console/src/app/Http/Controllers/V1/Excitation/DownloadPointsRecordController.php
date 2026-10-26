<?php

namespace App\Http\Controllers\V1\Excitation;

use App\Http\Controllers\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Traits\ExcelTrait;
use Illuminate\Http\Request;

/**
 * 积分记录导出
 */
class DownloadPointsRecordController extends BaseController
{

    use ExcelTrait;

    /**
     * 下载模板
     *
     * @return void
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index(Request $request)
    {
        set_time_limit(0);
        $limit         = 5000;
        $header        = ['导购id', '导购姓名', '所在部门', '关联规则名称', '积分维度',
                          '规则类目', '积分获取条件', '积分', '金币', '积分来源', '获取时间','财年财月'];
        $incentiveType = [
            1 => "邀约维度",
            2 => "任务维度",
            3 => '客户互动维度',
            4 => '培训维度',
            5 => '线上能力维度'
        ];
        $sourceType    = [
            1 => "规则统计",
            2 => "手动补发",
            3 => '奖章奖励',
        ];
        list($writer, $style) = $this->setExcelTitle($header, 'points-record');


        $page = 1;
        while (true) {
            $result = GoMicroScrm::pointsRecordList(array_merge($request->all(), [
                'enterprise_id' => $this->enterpriseId,
                'page'          => $page,
                'size'          => $limit
            ]));

            $dataArr = [];
            $listArr = $result['list'] ?? [];
            if (!empty($listArr)) {

                foreach ($listArr as $val) {
                    $dataArr[] = [
                        $val['member_user_id'],
                        $val['member_name'],
                        $val['department_name'],
                        $val['integral_rule_name'],
                        $incentiveType[$val['incentive_type']] ?? '',
                        $val['rule_category_name'],
                        $val['points_remark'],
                        round(floatval($val['points']),1),
                        round(floatval($val['coins']),1),
                        $sourceType[$val['source_type']] ?? '',
                        date('Y-m-d H:i:s', $val['created_at']),
                        $val['fy_time_str'],
                    ];
                }
                $this->setExcelData($writer, $style, $dataArr);
            } else {
                break;
            }
            $page++;
        }
        $this->closeExcel($writer);
    }


}