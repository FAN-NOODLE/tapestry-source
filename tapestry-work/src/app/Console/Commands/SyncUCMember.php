<?php

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\SyncOrgLog;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\EnterpriseAdmin;
use Illuminate\Console\Command;

class SyncUCMember extends Command
{

    /**
     * 命令
     * @var string
     */
    protected $signature = 'sync_uc_member';

    /**
     * 描述信息
     * @var string
     */
    protected $description = 'UC同步通讯录';

    /**
     * SyncDataLakeInfo constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 入口方法
     */
    public function handle()
    {
        $enterpriseLists = (new Enterprise())->lists();
        $enterpriseList = $enterpriseLists['data'];
        foreach ($enterpriseList as $v){
            $params = [];
            $enterpriseAdmin = (new EnterpriseAdmin())->detail([
                'enterprise_id' => $v['id'],
                'mobile'        => $v['mobile']
            ])['data'];

            $params['enterprise_id'] = $v['id'];
            $params['qy_domain'] = $v['domain'];
            $params['admin_id'] = $enterpriseAdmin['id'];
            $params['admin_name'] = $enterpriseAdmin['username'];

            $this->syncUCMember($params);
        }
        return true;

    }


    private function syncUCMember(array $params): bool
    {

        // 当天有同步未完成的禁止重复操作
        $count = SyncOrgLog::query()->where([
            ['enterprise_id', '=', $params['enterprise_id']],
            ['status', '<>', SyncOrgLog::STATUS_UC_OK],
        ])->count();
        if ($count) {
            $this->info("同步操作还未完成，请耐心等待！");
            return false;
        }
        // 调UC同步通讯录
        $ucResult = (new Member())->syncContact($params['qy_domain']);

        if (empty($ucResult['requestId'])) {
            $this->info("uc同步失败");
            return false;
        }
        // 同步操作记录
        unset($params['qy_domain']);
        $params['source_type'] = SyncOrgLog::SOURCE_TYPE_AUTO;
        SyncOrgLog::query()->create($params);

        return true;

    }

}