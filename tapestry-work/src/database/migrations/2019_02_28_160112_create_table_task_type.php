<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableTaskType extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('task_type')) {
            Schema::create('task_type', function (Blueprint $table) {
                $table->increments('id')->comment('任务类型id');
                $table->tinyInteger('serve_type')->default(0)->comment('任务类型偏向:0=crs+shop;1=shop;2=crs')->index();
                $table->string('name', 128)->default('')->comment('类型名称');
                $table->string('describe', 128)->default('')->comment('类型描述');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->unique('name');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `task_type` comment '任务类型表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('task_type');
    }
}
