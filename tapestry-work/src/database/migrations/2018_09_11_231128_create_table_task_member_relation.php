<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableTaskMemberRelation extends Migration
{
    /**
     * Run the migrations.
     * 导购任务表
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('task_member_relation')) {
            Schema::create('task_member_relation', function (Blueprint $table) {
                $table->increments('id')->comment('自增id');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业ID')->index();
                $table->unsignedInteger('task_id')->default(0)->comment('任务 ID')->indx();
                $table->char('mem_uid', 32)->default('')->comment('导购 ID')->index();
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `task_member_relation` comment '导购任务表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('task_member_relation');
    }
}
