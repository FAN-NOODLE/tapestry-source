<?php
/**
 * 营销计划任务进度记录表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateMarketingPlanTaskProgressTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'marketing_plan_task_progress';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('task_id')->default(0)->comment('营销计划任务表ID');
                $table->string('mem_uid', 150)->default('')->comment('员工uid');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('员工ID');
                $table->tinyInteger('finished_type')->default(1)->comment('完成方式(1、复制话术，2、转发给顾客，3、分享朋友圈)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('task_id', 'task_id_index');
            });

            DB::statement("alter table `{$tableName}` COMMENT '营销计划任务进度记录表' ");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //Schema::dropIfExists('marketing_plan');
    }
}
