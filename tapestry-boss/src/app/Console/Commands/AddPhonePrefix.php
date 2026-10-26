<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 8:01 PM
 */

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Models\EnterpriseAdmin;

class AddPhonePrefix extends Command
{
    /**
     * @var string
     */
    protected $signature = 'add_phone_prefix';
    /**
     * @var string
     */
    protected $description = '添加账号区号前缀';

    public function handle()
    {
        // 区号
        $prefix = '+86-';
        // 获取需要修改的数据
        $data = EnterpriseAdmin::query()->get(['account','id'])->toArray();

        // 数据处理
        foreach ($data as $key => $value) {
            // 判断是否已经有区号前缀
            if (strpos($value['account'], '+') === 0 && strpos($value['account'], '-')) {
                continue;
            }
            $value['account'] = $prefix . $value['account'];
            EnterpriseAdmin::query()->where('id', $value['id'])->update($value);
        }

        echo 'SUCCESS'.PHP_EOL;
    }
}
