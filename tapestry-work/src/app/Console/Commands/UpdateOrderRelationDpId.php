<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\StaffOrderDpRelation;
use App\Models\StaffOrderRelation;
use Illuminate\Console\Command;

class UpdateOrderRelationDpId extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-order-relation-dp-id';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新内购订单部门ID';

    public function handle()
    {

        $total = StaffOrderDpRelation::query()->count();
        if ($total) {
            echo '内购员工部门关系表已更新无需重复更新';
            return true;
        }

        // 获取订单关联关系列表
        $list = StaffOrderRelation::query()
            ->get(['id', 'order_no', 'enterprise_id', 'dp_id', 'mem_uid', 'created_at'])
            ->toArray();

        echo '内购员工数据共' . count($list) . '条';

        $data = [];
        $i    = 0;
        foreach ($list as $key => $item) {
            $data[] = [
                'enterprise_id'           => $item['enterprise_id'],
                'order_no'                => $item['order_no'],
                'mem_uid'                 => $item['mem_uid'],
                'dp_id'                   => $item['dp_id'],
                'staff_order_relation_id' => $item['id'],
                'created_at'              => $item['created_at']
            ];
            // 每次插入20条数据
            if ($i == 20) {
                StaffOrderDpRelation::query()->insert($data);
                $data = [];
                $i    = 0;
            } else {
                $i++;
            }
            echo '更新进度:' . ($key + 1) . '/' . count($list);
        }

        // 如果数据没有插入
        if ($data) {
            StaffOrderDpRelation::query()->insert($data);
        }

        echo '内购员工部门关系数据更新完成';

        return true;
    }

}
