<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Models\Activity;
use App\Models\ActivityReserveRecord;
use Exception;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Redis;
use Prettus\Repository\Exceptions\RepositoryException;

class AutoCIIEExpire extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:ciie:expire';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'CIIE吊牌预约活动过期处理';


    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return string
     *
     * @throws Exception
     */
    public function handle(): string
    {
        $id=4;
//        $enterprise_id=3363603667487424;
//        $activityRedisKey="CIIE:".$id;
//
//        $info = Activity::query()->where(['id' => $id, 'enterprise_id' => $enterprise_id])->first();
//        $activityInfo=Redis::connection()->set($activityRedisKey,json_encode($info,JSON_UNESCAPED_UNICODE));
//
//        return false;


        $reserve_time = time();
        // 查询预约活动
        ActivityReserveRecord::query()->where(["activity_id"=>$id])->where('reserve_time', '<', $reserve_time)->where(['is_check' => 0])->update(['is_check' => 2]);
        $list=ActivityReserveRecord::query()->where(["activity_id"=>$id])->where('reserve_time', '<', $reserve_time)->where(['is_check'=>2])->get();
        $userRecordKey="CIIE:record:{$id}";
        foreach ($list as $item){
            Redis::connection()->hset($userRecordKey,$item['openid'],json_encode($item));
        }
        return 'ok';
    }
}