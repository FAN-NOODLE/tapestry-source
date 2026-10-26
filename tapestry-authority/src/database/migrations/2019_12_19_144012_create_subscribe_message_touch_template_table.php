<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateSubscribeMessageTouchTemplateTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('subscribe_message_touch_template')) {
            Schema::create('subscribe_message_touch_template', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('主键id');
                $table->unsignedInteger('tid')->default(0)->index()->comment('微信模板标题id');
                $table->unsignedBigInteger('setting_id')->default(0)->index()->comment('配置id');
                $table->string('pri_tmpl_id', 191)->default('')->comment('微信个人模板库模板id');
                $table->string('touch_ids')->default('')->comment('触发事件id集合，json存储');
                $table->unsignedBigInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("alter table `subscribe_message_touch_template` comment '订阅消息-事件模板表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('subscribe_message_touch_template');
    }
}
