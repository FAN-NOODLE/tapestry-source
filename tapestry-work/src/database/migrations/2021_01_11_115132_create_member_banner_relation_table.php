<?php
/**
 * 导购小程序太阳码海报表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateMemberBannerRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'member_banner_relation';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('task_id')->default(0)->comment('营销计划任务表ID');
                $table->string('mem_uid', 150)->default('')->comment('员工uid');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('员工ID');
                $table->string('background_url', 255)->default('')->comment('背景海报路径');
                $table->string('suncode_mediaid', 255)->default('')->comment('生成的太阳码海报mediaid');
                $table->string('suncode_url', 255)->default('')->comment('太阳码路径');
                $table->text('suncode_param')->default('')->comment('太阳码json参数');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('task_id', 'task_id_index');
            });

            DB::statement("alter table `{$tableName}` COMMENT '导购小程序太阳码海报表' ");
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
