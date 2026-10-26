<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/30 0030
 * Time: 10:14
 */

namespace App\Console\Commands;


use App\Jobs\ExternalContactUpdateCustomerStatisticJob;
use App\Library\DataLake\UserRequestService;
use App\Library\Uc\External;
use App\Library\Uc\Member;
use App\Models\MemberUserRelation;
use App\Models\StoreWorker;
use App\Models\User;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class AutoExternalMissing extends Command
{

    /**
     * 命令
     * @var string
     */
    protected $signature = 'auto_external_missing';

    /**
     * 描述信息
     * @var string
     */
    protected $description = '同步UC缺失的外部联系人';


    /**
     * 入口方法
     */
    public function handle()
    {
        $data       = (new Enterprise())->listsOfficial()['data'];
        foreach ($data as $val) {
            $res = (new External())->syncMissing($val['domain']);
            print_r($res);
        }
        return true;
    }



}