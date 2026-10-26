<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRolePermissionTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('role_permission')) {
            Schema::create('role_permission', function (Blueprint $table) {
                $table->increments('id')->comment('id-自增');
                $table->unsignedInteger('role_id')->default(0)->index()->comment('角色id');
                $table->unsignedInteger('bap_id')->default(0)->index()->comment('权限id');
                $table->unsignedTinyInteger('type')->default(1)->comment('权限类型：1=基础，2=应用,默认基础');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `role_permission` COMMENT '角色权限关系表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('role_permission');
    }
}
