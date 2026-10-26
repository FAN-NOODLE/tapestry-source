<?php
/**
 * 营销计划任务表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateMarketingPlanTaskTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'marketing_plan_task';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('plan_id')->default(0)->comment('营销计划表ID');
                $table->string('task_title', 255)->default('')->comment('任务名称');
                $table->unsignedInteger('start_time')->default(0)->comment('任务开始时间');
                $table->unsignedInteger('end_time')->default(0)->comment('任务结束时间');
                $table->text('content')->default('')->comment('营销话术');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('plan_id', 'plan_id_index');
            });

            DB::statement("alter table `{$tableName}` COMMENT '营销计划任务表' ");
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
