<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableTaskMemberRecord extends Migration
{
    /**
     * Run the migrations.
     * 导购下员工任务记录表
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('task_member_record')) {
            Schema::create('task_member_record', function (Blueprint $table) {
                $table->increments('id')->comment('自增id');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业ID')->index();
                $table->unsignedInteger('task_id')->default(0)->comment('任务 ID');
                $table->char('mem_uid', 32)->default('')->comment('导购 ID')->index();
                $table->string('customer_id', 50)->default('')->comment('客户ID');
                $table->tinyInteger('task_type')->default(0)->comment('任务类型：(1:销售业绩，2:客户管理，3:生日祝福，4:礼券到期提醒，5:会员招募)');
                $table->tinyInteger('display_state')->default(1)->comment('是否显示：(1:显示,2:隐藏)');
                $table->tinyInteger('status')->default(1)->comment('完成状态：(1:未完成,2:已完成)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `task_member_record` comment '导购下员工任务记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('task_member_record');
    }
}
