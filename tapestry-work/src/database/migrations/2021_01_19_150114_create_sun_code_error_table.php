<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateSunCodeErrorTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'sun_code_error';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('task_id')->default(0)->comment('营销计划任务表ID');
                $table->text('sun_code_param')->default('')->comment('太阳码json参数');
                $table->text('error')->default('')->comment('错误信息');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `{$tableName}` COMMENT '太阳码生成错误表' ");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('sun_code_error');
    }
}
