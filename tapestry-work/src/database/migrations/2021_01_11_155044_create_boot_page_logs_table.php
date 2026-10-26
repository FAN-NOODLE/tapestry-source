<?php
/**
 * 引导页浏览记录表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateBootPageLogsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'boot_page_logs';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('version_code', 255)->default('')->comment('版本号');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('员工ID');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `{$tableName}` COMMENT '引导页浏览记录表' ");
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
