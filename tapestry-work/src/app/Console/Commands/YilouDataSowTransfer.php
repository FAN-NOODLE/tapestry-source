<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月20日17:9:30
 */

namespace App\Console\Commands;

use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;

class YilouDataSowTransfer extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'yilou-data-sow-transfer';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '一楼数据迁移清洗';

    //分页大小
    protected $size = 1000;
    // 最后一页
    protected $lastPage = 0;
    //页码
    protected $page = 1;

    public function handle()
    {

        $enterpriseId = env("COACH_ENTERPRISE_ID", 3360610327843968);// 蔻驰企业数据

        $lab = true;

        while ($lab) {
            $data = DB::connection('yilou_data')
                ->table("coach_customer_remind")
                ->paginate($this->size, ['*'], 'page', $this->page);

            $this->lastPage = $data->lastPage();
            if ($this->lastPage == $this->page) {
                $lab = false;
            }
            $this->page++;


            $unionIdArr    = [];
            $memberCodeArr = [];
            foreach ($data->items() as $val) {
                $unionIdArr[]    = $val->UnionID;
                $memberCodeArr[] = $val->EmpnoCode;
            }

            $externalUser = WxUserCorpExternal::query()->setConn($enterpriseId)
                ->whereIn('unionid', $unionIdArr)
                ->get(['id', 'unionid', 'nick'])
                ->keyBy('unionid')->toArray();

            $memberArr = Members::query()->where('enterprise_id', $enterpriseId)
                ->whereIn('userid', $memberCodeArr)
                ->get(['id', 'name', 'userid'])
                ->keyBy('userid')
                ->toArray();


            $addArr = [];

            foreach ($data as $val) {
                $val     = (array)$val;
                $tmpData = [
                    'id'                    => $val['Id'],
                    'enterprise_id'         => $enterpriseId,
                    'work_external_user_id' => isset($externalUser[$val['UnionID']]) ? $externalUser[$val['UnionID']]['id'] : 0,
                    'member_id'             => isset($memberArr[$val['EmpnoCode']]) ? $memberArr[$val['EmpnoCode']]['id'] : 0,
                    'member_user_id'        => $val['EmpnoCode'],
                    'content'               => $val['RemindContent'],
                    'status'                => $val['Status'] == 5 ? 2 : 1,
                    'remind_time'           => strtotime($val['RemindDate']),
                    'created_at'            => strtotime($val['CreateTime']),
                    'updated_at'            => strtotime($val['UpdateTime'])
                ];

                $addArr[] = $tmpData;
            }

            Config::set('database.connections.micro_scrm.host', '10.1.1.8');
            Config::set('database.connections.micro_scrm.port', 3306);
            Config::set('database.connections.micro_scrm.username', 'root');
            Config::set('database.connections.micro_scrm.password', '6F*Ej^Eqz8fmUj1T');
            Config::set('database.connections.micro_scrm.database', 'micro_scrm');
            DB::connection('micro_scrm')->reconnect();
            $data = DB::connection('micro_scrm')->table("sow_task")->insert($addArr);
            print_r($data);
        }
    }
}