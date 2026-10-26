<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/8/4
 * Time: 10:01 下午
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\MarketingPlan;
use Maatwebsite\Excel\Facades\Excel;


class ExportController extends BaseController
{

    /**
     * 设置导出列宽
     *
     * @var array
     */
    protected $colWidth = [
        'A' => 20,
        'B' => 20,
        'C' => 20,
        'D' => 15,
    ];

    /**
     * 导出导购列表
     */
    public function index()
    {
        $input                  = $this->request->only(['plan_id', 'task_id']);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['file_name']     = '未完成导购列表';
        $data                   = (new MarketingPlan())->notFinishedList($input);
        return $this->export($data);
    }

    private function export($lists)
    {
        $cellData[] = [
            '导购姓名',
            '导购编号',
            '电话',
            '部门名'
        ];
        if ($lists) {
            foreach ($lists['not_finished'] ?? [] as $value) {
                $cellData[] = [
                    (string)$value['name'],
                    (string)$value['worker_no'],
                    $value['mobile'],
                    $value['store_name']
                ];
            }
        }
        $name = '未完成人员名单' . date('YmdHi', time());
        return Excel::create($name, function ($excel) use ($cellData) {
            $excel->sheet('not_finished', function ($sheet) use ($cellData) {
                //$sheet->setWidth($this->colWidth)->rows($cellData);
                $sheet->setWidth($this->colWidth)->rows($cellData);
            });
        })->export('xlsx');
    }
}