<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 9:51
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: StatisticsController.php
 */

namespace App\Http\Controllers\Brand\CommunityStatistics;


use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Work\DataCompass;
use Maatwebsite\Excel\Facades\Excel;

class ExportController extends BaseController
{
    /**
     * 设置导出列宽
     * @var array
     */
    protected $colWidth = [
        'A' => 20,
        'B' => 15,
        'C' => 15,
        'D' => 15,
        'E' => 15,
        'F' => 15,
        'G' => 15,
        'H' => 15,
//        'I' => 15,
//        'J' => 15,
    ];

    protected $baseDataKey = [
        '累计群成员数',
        '新入群人数',
        '退群人数',
        '客户入群率',
        '客户流失率',
        '加好友人数',
        '加好友率',
//        '微信绑定会员数',
//        '微信绑定会员率',
        '群活动数',
        '群活动参与人数',
        '未参加群活动人数',
        '未参加群活动率',
    ];

    /**
     * 导出列表类型（1-按社群，2-按导购）
     */
    const BY_GROUP_EXPORT = 1;
    const BY_MEMBER_EXPORT = 2;

    /**
     * 导出社群统计数据明细列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:51
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function index(Request $request)
    {
        $params                  = $request->input();
        $params['type']          = $request->input('type', self::BY_GROUP_EXPORT);
        $params['start_time']    = $request->input('start_time');
        $params['end_time']      = $request->input('end_time');
        $params['worker_ids']    = $request->input('worker_ids', []);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain']     = $this->qyDomain;
        $params['is_page']       = false;

        if(!in_array($params['type'], [self::BY_GROUP_EXPORT, self::BY_MEMBER_EXPORT])){
            notice(2300002);
        }

        if($params['type'] == self::BY_GROUP_EXPORT){
            $statistics = (new DataCompass())->getCommunityStatisticListByGroup($params);
        }else{
            $statistics = (new DataCompass())->getCommunityStatisticListByMember($params);
        }

        return $this->export($statistics, $params['type']);
    }

    protected function export($lists, $type)
    {
        if($type == self::BY_GROUP_EXPORT){
            $cellData[] = array_merge([
                '社群名称'
            ], $this->baseDataKey);
        }else{
            $cellData[] = array_merge([
                '导购ID',
                '导购头像',
                '导购昵称',
                '导购所在门店'
            ], $this->baseDataKey);
        }

        if ($lists) {
            foreach ($lists ?? [] as $value) {
                $baseData = [
                    (int)$value['group_member_count'],
                    (int)$value['new_join_group_count'],
                    (int)$value['leave_group_count'],
                    (string)$value['join_group_rate'],
                    (string)$value['leave_group_rate'],
                    (int)$value['join_friend_count'],
                    (string)$value['join_friend_rate'],
//                    (int)$value['wechat_binding_vip_count'],
//                    (string)$value['wechat_binding_vip_rate'],
                    (int)$value['group_activity_count'],
                    (int)$value['join_group_activity_count'],
                    (int)$value['not_join_group_activity_count'],
                    (string)$value['not_join_group_activity_rate'],
                ];

                if($type == self::BY_GROUP_EXPORT){
                    $cellData[] = array_merge([
                        (string)$value['name'],
                    ], $baseData);
                }else{
                    $cellData[] = array_merge([
                        (string)$value['worker_no'],
                        (string)$value['avatar'],
                        (string)$value['nickname'],
                        (string)$value['department_name'],
                    ], $baseData);
                }

            }
        }

        $name = '社群数据明细' . date('YmdHi', time());
        return Excel::create($name, function ($excel) use ($cellData) {
            $excel->sheet('社群数据明细', function ($sheet) use ($cellData) {
                $sheet->setWidth($this->colWidth)->rows($cellData);
            });
        })->export('xlsx');
    }
}