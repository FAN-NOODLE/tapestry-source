<?php

namespace App\Http\Controllers\V1\Excitation;

use App\Http\Controllers\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Traits\ExcelTrait;

/**
 * 奖章领取记录
 */
class MedalRecordController extends BaseController
{

    use ExcelTrait;

    public function download()
    {
        set_time_limit(0);

        $params = $this->request->all();
        $isCoach=in_array($this->enterpriseId,[3360610327843968,3363603667487424]);
        $limit  = 2000;
        $header = ['导购ID', '导购姓名', '所在部门', '关联规则名称', '奖章名称', '奖章来源', '奖章规则维度', '奖章规则类目', '奖章获取标准', '实际达成情况', '实际业务数值', '积分', '获取时间'];
        if ($isCoach){
            $header = ['导购ID', '导购姓名', '所在部门', '关联规则名称', '奖章名称', '奖章来源', '奖章规则维度', '奖章规则类目', '奖章获取标准', '实际达成情况', '实际业务数值', '积分', '金币', '实物奖品', '获取时间'];
        }

        $sourceType = [1 => "手动补发", 2 => "规则统计"];
        list($writer, $style) = $this->setExcelTitle($header, 'points-record');

        $page = 1;
        while (true) {
            $requestParams = array_merge($params, [
                'enterprise_id' => $this->enterpriseId,
                'page'          => $page,
                'size'          => $limit
            ]);
            $result        = GoMicroScrm::medalRecordList($requestParams);
            if (!isset($result['list'])) {
                break;
            }

            $dataArr = [];
            foreach ($result['list'] as $val) {
                $data=[
                    $val['member_user_id'],
                    $val['member_name'],
                    $val['department_name'],
                    $val['rule_name'],
                    $val['medal_name'],
                    $sourceType[$val['source']] ?? '',
                    $val['rule_dimension_name'],
                    $val['rule_category_name'],
                    $val['description'] . $val['invite_top'],
                    $val['ranking'],
                    $val['ranking_value'],
                    $val['integral'],
                ];
                if ($isCoach){
                    $data=array_merge($data,[
                        floatval($val['coins']),
                        $val['item_name'],
                        date('Y-m-d H:i:s', $val['created_at'])
                    ]);
                }else{
                    $data=array_merge($data,[date('Y-m-d H:i:s', $val['created_at'])]);
                }
                $dataArr[] = $data;
            }
            $this->setExcelData($writer, $style, $dataArr);

            if (count($result['list']) < $limit) {
                break;
            }
            $page++;
        }
        $this->closeExcel($writer);
    }
}