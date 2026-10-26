<?php
/**
 * Created by PhpStorm.
 * User: heshupeng
 * Date: 2018/12/5
 * Time: 15:58
 */

namespace App\Console\Commands\Message;

use App\Models\EnterpriseThirdPartAccount;
use App\Services\SubscribeMessageService;
use Illuminate\Console\Command;

class AutoSyncTemplates extends Command
{

    /**
     * @var string
     */
    protected $signature = 'auto:sync_templates {--enterprise_id=}';

    /**
     * @var string
     */
    protected $description = '更新小程序模板消息';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * @return bool
     */
    public function handle()
    {
        $enterpriseId = $this->option('enterprise_id');
        $list = $this->getMiniList(['version' => 'gw-shop', 'enterprise_id' => $enterpriseId]);

        // 循环小程序数据
        foreach ($list as $value) {
            $startTime = microtime(1);
            $this->info('企业' . $value['enterprise_id'] . '开始处理,时间戳ms：' . $startTime);
            $res = (new SubscribeMessageService($value['enterprise_id']))->sync();
            $this->info('企业' . $value['enterprise_id'] . '处理结果：' . json_encode($res, 320));
            $this->info('企业' . $value['enterprise_id'] . '处理完成,耗时：' . (microtime(1) - $startTime));
        }

        return true;
    }

    /**
     * @param array $params
     * @return array
     */
    public function getMiniList(array $params = [])
    {
        return EnterpriseThirdPartAccount::query()
            ->when($params['enterprise_id'], function ($query) use ($params) {
                $query->where('enterprise_id', $params['enterprise_id']);
            })
            ->where('enterprise_third_part_account.app_id', '<>', '')
            ->where('enterprise_third_part_account.app_secret', '<>', '')
            ->where('enterprise_third_part_account.status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('enterprise_third_part_account.platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('wx_category_ids', '<>', '')
            ->join('enterprise', 'enterprise.id', '=', 'enterprise_third_part_account.enterprise_id')
            ->when('shop' == ($params['version'] ?? ''), function ($query) {
                // 老版
                return $query->where('enterprise.version', 1);
            })
            ->when('gw-shop' == ($params['version'] ?? ''), function ($query) {
                // 老版
                return $query->where('enterprise.version', '<>', 1);
            })
            ->orderBy('enterprise_third_part_account.enterprise_id', 'asc')
            ->get([
                'enterprise_third_part_account.enterprise_id',
                'enterprise_third_part_account.app_id',
                'enterprise_third_part_account.app_secret',
                'enterprise_third_part_account.wx_category',
                'enterprise_third_part_account.wx_category_ids',
                'enterprise.version'
            ])
            ->toArray();
    }
}
