<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRoleTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('role', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary();
            $table->unsignedBigInteger('pid')->default(0)->comment('父ID');
            $table->unsignedBigInteger('admin_id')->default(0)->comment('管理员ID');
            $table->string('mark', 16)->default('')->comment('标识');
            $table->string('name', 64)->default('')->comment('角色名称');
            $table->unsignedInteger('weight')->default(0)->comment('权重');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');

            $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
        });

        \Illuminate\Support\Facades\DB::statement("alter table `role` comment '角色表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('role');
    }
}
