<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateSubscribeMessageRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('subscribe_message_record')) {
            Schema::create('subscribe_message_record', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('主键id');
                $table->unsignedBigInteger('user_id')->default(0)->index()->comment('用户id');
                $table->unsignedInteger('tid')->default(0)->index()->comment('微信模板标题id');
                $table->string('pri_tmpl_id', 200)->default('')->comment('微信个人模板库模板id');
                $table->unsignedTinyInteger('type')->default(0)->comment('类型（1：允许一次，2：总是允许）');
                $table->unsignedInteger('send_num')->default(0)->comment('发送次数');
                $table->unsignedTinyInteger('object_type')->default(0)->index()->comment('对象类型（1：视频直播）');
                $table->unsignedBigInteger('object_id')->default(0)->comment('对象类型对应的对象id');
                $table->unsignedBigInteger('touch_id')->default(0)->comment('触发事件id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("alter table `subscribe_message_record` comment '订阅消息-订阅记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('subscribe_message_record');
    }
}
