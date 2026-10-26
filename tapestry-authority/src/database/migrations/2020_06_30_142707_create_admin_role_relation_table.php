<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateAdminRoleRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('admin_role_relation', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary();
            $table->unsignedBigInteger('admin_id')->default(0)->comment('管理员ID');
            $table->unsignedBigInteger('role_id')->default(0)->comment('角色ID');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');

            $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
        });

        \Illuminate\Support\Facades\DB::statement("alter table `admin_role_relation` comment '管理员角色关系表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('admin_role_relation');
    }
}
