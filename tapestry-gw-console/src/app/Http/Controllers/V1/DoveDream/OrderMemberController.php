<?php
/**
 * Created by PhpStorm.
 * User: qihongchuan
 * Date: 2020/02/21
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\DoveDream;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DoveDream;
use Maatwebsite\Excel\Facades\Excel;

class OrderMemberController extends BaseController
{
    // 参加者性别 ： 男
    const MEMBER_GENDER_MALE = 1;
    // 参加者性别 ： 女
    const MEMBER_GENDER_FEMALE = 2;

    /**
     * 参观者详情
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params = $request->only(['id']);
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getOrderMemberDetail($params);
        return success($result);
    }

    /**
     * 参观者编辑
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function save(Request $request)
    {
        $params = $request->only(['id','name','age','gender','mobile','id_card']);
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->postOrderMemberSave($params);
        return success($result);
    }

    /**
     * 参观者列表
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function list(Request $request)
    {
        $params = $request->only(['type','keyword','start_time','end_time','status','sort','page','size']);
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getOrderMemberList($params);
        return success($result);
    }

    /**
     * 参观者导出
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function export(Request $request)
    {
        $params = $request->only(['type','keyword','start_time','end_time','status']);
        $params['enterprise_id'] = $this->enterpriseId;

        // 查询数据
        $memberList = (new DoveDream)->getOrderMemberExport($params)['data'];

        // 列宽
        $colWidth = ['A' => 25,'B' => 20,'C' => 15,'D' => 20,'E' => 15,'F' => 15,'G' => 15,'H' => 40,'I' => 13,'J' => 10,'K' => 18];

        // 导出数据
        $cellData[] = ['订单编号','下单客户','下单时间','姓名','年龄','性别','手机号','身份证号','参加时间','核销码','状态'];
        // todo : 字段展示/导出格式调整，如金额等
        foreach ($memberList as $member) {
            $cellData[] = [
                "\t".$member['order_no'] ?? '',
                $member['nick_name'] ?? '',
                empty($member['created_at']) ? '' : date('Y-m-d H:i:s', (int)$member['created_at']),
                $member['name'] ?? '',
                $member['age'] ?? '',
                $member['gender'] == $this::MEMBER_GENDER_MALE ? '男' : '女',
                "\t".$member['mobile'] ?? '',
                "\t".$member['id_card'] ?? '',
                $member['visit_term'] ?? '',
                "\t".$member['write_off_no'] ?? '',
                $member['write_off_at'] > 0 ? '已核销' : '待核销',
            ];
        }
        // TODO : 导出文件名乱码问题
        $exportFileName = 'member-list-export-' . date('YmdHi', time());

        return Excel::create($exportFileName, function ($excel) use ($cellData, $colWidth) {
            $excel->sheet('score', function ($sheet) use ($cellData, $colWidth) {
                $sheet->setWidth($colWidth)->rows($cellData);

                $endLine = count($cellData);
                // 设置表头样式
                $sheet->cells('A1:X1', function ($cells) {
                    $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                })->cells('A2:X' . $endLine, function ($cells) { // 设置内容样式
                    $cells->setAlignment('center')->setValignment('center');
                });
            });
        })->export('xlsx');
    }

    /**
     * 参观者核销
     * @author zhaoqi
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function write(Request $request)
    {
        $input = $request->only([
            'write',
        ]);

        $rules = [
            'write' => 'required'
        ];

        $message = [
            'write.required'  => '核销id必须',
        ];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_id'] = $this->enterpriseAdminId;
        $result = (new DoveDream)->write($input);
        return success($result);
    }
}

