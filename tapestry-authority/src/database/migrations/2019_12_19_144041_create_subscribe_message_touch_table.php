<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateSubscribeMessageTouchTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('subscribe_message_touch')) {
            Schema::create('subscribe_message_touch', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('主键id');
                $table->string('name', 200)->default('')->comment('名称');
                $table->string('comment')->default('')->comment('备注');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("alter table `subscribe_message_touch` comment '订阅消息-触发事件表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('subscribe_message_touch');
    }
}
