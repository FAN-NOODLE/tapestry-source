<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月20日17:9:30
 */

namespace App\Console\Commands;

use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Models\StoreDepartment;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;

class YilouDataTransfer extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'yilou-data-transfer';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '一楼数据迁移清洗';


    public function handle()
    {
        $data          = DB::connection('yilou_data')
            ->select('select * from coach_customer_appointments limit 10');
        $unionidArr    = [];
        $memberCodeArr = [];
        foreach ($data as $val) {
            $unionidArr[]    = $val->UnionID;
            $memberCodeArr[] = $val->EmpnoCode;
        }
        $st = StoreDepartment::query()->where('enterprise_id', env('COACH_ENTERPRISE_ID'))
            ->get(['wechat_department_id', 'third_store_code', 'store_name', 'address', 'store_type'])
            ->keyBy('third_store_code')->toArray();

        $external   = WxUserCorpExternal::query()->setConn(env('COACH_ENTERPRISE_ID'))
            ->whereIn('unionid', $unionidArr)
            ->get(['unionid', 'nick'])
            ->keyBy('unionid')->toArray();
        $addArr     = [];
        $serviceArr = [
            12 => 3389002357043584,
            13 => 3389002454845824,
            14 => 3389002559508864
        ];
        $memberArr  = Members::query()->where('enterprise_id', env('COACH_ENTERPRISE_ID'))
            ->whereIn('userid', $memberCodeArr)
            ->get(['id', 'name', 'userid'])
            ->keyBy('userid')->toArray();
        foreach ($data as $val) {
            $val      = (array)$val;
            $addArr[] = [
                'id'                  => $val['Id'],
                'reserve_number'      => $val['AppointId'],
                'reserve_id'          => $serviceArr[$val['ServiceID']] ?? $val['ServiceID'],
                'reserve_time'        => strtotime($val['AppointmentDate']),
                'store_id'            => isset($st[$val['StoreCode']]) ? $st[$val['StoreCode']]['wechat_department_id'] : 0,
                'store_name'          => isset($st[$val['StoreCode']]) ? $st[$val['StoreCode']]['store_name'] : $val['StoreCode'],
                'store_address'       => isset($st[$val['StoreCode']]) ? $st[$val['StoreCode']]['address'] : '',
                'status'              => $val['Status'] == 2 ? 1 : ($val['Status'] == 1 ? 4 : 3),//yilou 1 默认 2 匹配订单成功 3 删除
                'user_name'           => isset($external[$val['UnionID']]) ? $external[$val['UnionID']]['nick'] : '',
                'user_unionid'        => $val['UnionID'],
                'user_external_id'    => $val['ExtUserId'],
                'member_id'           => isset($memberArr[$val['EmpnoCode']]) ? $memberArr[$val['EmpnoCode']]['id'] : 0,
                'member_name'         => isset($memberArr[$val['EmpnoCode']]) ? $memberArr[$val['EmpnoCode']]['name'] : $val['EmpnoCode'],
                'is_open'             => $val['CustomerAccept'],
                'is_commit'           => $val['CustomerAccept'],
                'invite_num'          => 1,
                'member_user_id'      => $val['EmpnoCode'],
                'reserve_order_sn'    => $val['OrderID'],
                'reserve_order_price' => intval($val['OrderAccount'] * 100),
                'reserve_order_time'  => strtotime($val['OrderDate']),
                'is_pay'              => $val['OrderID'] ? 1 : 2,
                'is_confirm'          => 2,
                'order_origin'        => $val['AppointmentType'],
                'created_at'          => strtotime($val['CreateTime']),
                'updated_at'          => strtotime($val['CreateTime'])
            ];
        }
        Config::set('database.connections.micro_scrm_enterprise.host', '10.2.1.13');
        Config::set('database.connections.micro_scrm_enterprise.port', 3306);
        Config::set('database.connections.micro_scrm_enterprise.username', 'root');
        Config::set('database.connections.micro_scrm_enterprise.password', 'GnKASV0IArjht^uM');
        Config::set('database.connections.micro_scrm_enterprise.database', 'micro_scrm_enterprise_3363603667487424');
        DB::connection('micro_scrm_enterprise')->reconnect();
//        $info = DB::connection('micro_scrm_enterprise')->table('reserve_order')->first();
//        dd($info);
        DB::connection('micro_scrm_enterprise')->table('reserve_order')->insert($addArr);

    }
}
