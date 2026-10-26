<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableTask extends Migration
{
    /**
     * Run the migrations.
     * 任务主表
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('task')) {
            Schema::create('task', function (Blueprint $table) {
                $table->increments('id')->comment('自增id');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业ID')->index();
                $table->string('task_name',191)->default('')->comment('任务名称');
                $table->tinyInteger('task_type')->default(0)->comment('任务类型：(1:销售业绩，2:客户管理，3:生日祝福，4:礼券到期提醒，5:会员招募)');
                $table->Integer('start_time')->default(0)->comment('开始时间');
                $table->Integer('end_time')->default(0)->comment('结束时间');
                $table->unsignedInteger('task_cycle')->default(1)->comment('任务周期：(1:一天，2:一周，3:一个月，4:三个月)');
                $table->unsignedInteger('is_loop')->default(2)->comment('是否循环：(1:是，2=否)');
                $table->Integer('task_index')->default(0)->comment('任务指标');
                $table->unsignedInteger('is_start_task')->default(1)->comment('任务开始之前提醒：(1:是，2:否)');
                $table->unsignedInteger('is_end_task')->default(1)->comment('任务结束之前提醒：(1:是，2:否)');
                $table->Integer('task_start_time')->default(0)->comment('开始提醒时间');
                $table->Integer('task_end_time')->default(0)->comment('结束提醒时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `task` comment '任务主表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('task');
    }
}
