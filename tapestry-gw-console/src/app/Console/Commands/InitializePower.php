<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2018/12/4
 * Time: 19:40
 */

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class InitializePower extends Command
{
    /**
     * @var string
     */
    protected $signature = 'initialize:initialize_power';

    protected $description = '初始化权限';

    protected $insert_data = [];

    public function __construct()
    {
        parent::__construct();
    }

    public function handle()
    {

        $data = [
            [
                'name'       => '首页',
                'alias'      => 'app.home',
                'path'       => 'home',
                'parent'     => '',
                'serve_type' => 0,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [],
            ],
            [
                'name'       => '首页',
                'alias'      => 'app.home',
                'path'       => 'home',
                'parent'     => '',
                'serve_type' => 1,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [],
            ],
            [
                'name'       => '首页',
                'alias'      => 'app.home',
                'path'       => 'home',
                'parent'     => '',
                'serve_type' => 2,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [],
            ],
            [
                'name'       => '任务中心',
                'alias'      => 'app.taskCenter',
                'path'       => 'taskCenter',
                'parent'     => '',
                'serve_type' => 2,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [
                    [
                        'name'       => '任务列表',
                        'alias'      => 'app.taskCenter.taskList',
                        'path'       => 'taskList',
                        'parent'     => 'taskCenter',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '任务类型',
                        'alias'      => 'app.taskCenter.taskType',
                        'path'       => 'taskType',
                        'parent'     => 'taskCenter',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                ],
            ],
            [
                'name'       => '商城管理',
                'alias'      => 'app.mallManagement',
                'path'       => 'mallManagement',
                'parent'     => '',
                'serve_type' => 1,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [
                    [
                        'name'       => '商品管理',
                        'alias'      => 'app.mall.goodsManagement',
                        'path'       => 'goodsManagement',
                        'parent'     => 'mallManagement',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '订单管理',
                        'alias'      => 'app.mall.orderManagement',
                        'path'       => 'orderManagement',
                        'parent'     => 'mallManagement',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '佣金规则',
                        'alias'      => 'app.mall.brokerageDetail',
                        'path'       => 'brokerageDetail',
                        'parent'     => 'mallManagement',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '页面管理',
                        'alias'      => 'app.mall.pageManagement',
                        'path'       => 'pageManagement',
                        'parent'     => 'mallManagement',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                ],
            ],
            [
                'name'       => '设置',
                'alias'      => 'app.setting',
                'path'       => 'setting',
                'parent'     => '',
                'serve_type' => 0,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [
                    [
                        'name'       => '微信授权',
                        'alias'      => 'app.setting.authorize',
                        'path'       => 'authorize',
                        'parent'     => 'setting',
                        'serve_type' => 0,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '修改密码',
                        'alias'      => 'app.setting.password',
                        'path'       => 'password',
                        'parent'     => 'setting',
                        'serve_type' => 0,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '支付配置',
                        'alias'      => 'app.setting.wxpay',
                        'path'       => 'wxpay',
                        'parent'     => 'setting',
                        'serve_type' => 0,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '客户联系配置',
                        'alias'      => 'app.setting.contact',
                        'path'       => 'contact',
                        'parent'     => 'setting',
                        'serve_type' => 0,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                ],
            ],
            [
                'name'       => '设置',
                'alias'      => 'app.setting',
                'path'       => 'setting',
                'parent'     => '',
                'serve_type' => 1,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [
                    [
                        'name'       => '微信授权',
                        'alias'      => 'app.setting.authorize',
                        'path'       => 'authorize',
                        'parent'     => 'setting',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '修改密码',
                        'alias'      => 'app.setting.password',
                        'path'       => 'password',
                        'parent'     => 'setting',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '支付配置',
                        'alias'      => 'app.setting.wxpay',
                        'path'       => 'wxpay',
                        'parent'     => 'setting',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '客户联系配置',
                        'alias'      => 'app.setting.contact',
                        'path'       => 'contact',
                        'parent'     => 'setting',
                        'serve_type' => 1,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                ],
            ],
            [
                'name'       => '设置',
                'alias'      => 'app.setting',
                'path'       => 'setting',
                'parent'     => '',
                'serve_type' => 2,
                'created_at' => time(),
                'updated_at' => time(),
                'deleted_at' => time(),
                'child'      => [
                    [
                        'name'       => '微信授权',
                        'alias'      => 'app.setting.authorize',
                        'path'       => 'authorize',
                        'parent'     => 'setting',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '修改密码',
                        'alias'      => 'app.setting.password',
                        'path'       => 'password',
                        'parent'     => 'setting',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '支付配置',
                        'alias'      => 'app.setting.wxpay',
                        'path'       => 'wxpay',
                        'parent'     => 'setting',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                    [
                        'name'       => '客户联系配置',
                        'alias'      => 'app.setting.contact',
                        'path'       => 'contact',
                        'parent'     => 'setting',
                        'serve_type' => 2,
                        'created_at' => time(),
                        'updated_at' => time(),
                        'deleted_at' => time(),
                        'child'      => [],
                    ],
                ],
            ],
        ];


        $data = $this->format_data($data);

        DB::table('base_power')->insert($data);
    }


    public function format_data($data)
    {

        foreach ($data as $v) {

            if (empty($v['child'])) {

                unset($v['child']);
                $this->insert_data[] = $v;

            } else {

                $arr = $v['child'];
                $this->format_data($arr);
                unset($v['child']);
                $this->insert_data[] = $v;
            }
        }

        return $this->insert_data;
    }
}
