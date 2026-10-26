<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseAdminRole extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_admin_role')) {
            Schema::create('enterprise_admin_role', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->default(0)->comment('管理员角色关系主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('enterprise_admin_id')->default(0)->comment('企业账号表id')->index();
                $table->unsignedBigInteger('role_id')->default(0)->comment('角色表主键ID')->index();
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("alter table `enterprise_admin_role` comment '管理员角色关联表'");
        }

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_admin_role');
    }
}
