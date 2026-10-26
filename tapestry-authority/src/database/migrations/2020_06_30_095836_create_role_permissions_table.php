<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRolePermissionsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('role_permissions', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary();
            $table->unsignedBigInteger('role_id')->default(0)->comment('角色ID');
            $table->unsignedBigInteger('permissions_id')->default(0)->comment('权限ID');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');

            $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
        });

        \Illuminate\Support\Facades\DB::statement("alter table `role_permissions` comment '角色权限表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('role_permissions');
    }
}
