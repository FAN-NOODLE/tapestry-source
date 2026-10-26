<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseActivityMessage extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_activity_message')) {
            Schema::create('enterprise_activity_message', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary();
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
                $table->unsignedInteger('message_scene_id')->default(0)->comment('订阅消息策略id');
                $table->string('title')->default('')->comment('活动标题');
                $table->unsignedInteger('send_time')->default(0)->comment('发送时间');
                $table->unsignedInteger('activity_start_time')->default(0)->comment('活动开始时间');
                $table->unsignedTinyInteger('page_type')->default(0)->comment('跳转类型 1 商品 2 活动页 3 直播列表 4 直播详情 5 活动 6 页面路径');
                $table->string('page_url')->default('')->comment('跳转地址');
                $table->string('page_query')->default('')->comment('跳转地址参数');
                $table->unsignedTinyInteger('send_status')->default(0)->comment('发送状态 0 未发送 1 已发送');
                $table->unsignedTinyInteger('enable_status')->default(0)->comment('启用状态 0 正常 1 已停用');
                $table->unsignedInteger('expected_num')->default(0)->comment('预计发送条数');
                $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `enterprise_activity_message` comment '活动提醒配置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('enterprise_activity_message')) {
            Schema::dropIfExists('enterprise_activity_message');
        }
    }
}
