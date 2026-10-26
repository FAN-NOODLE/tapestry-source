<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseRoleTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_role')) {
            Schema::create('enterprise_role', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->default(0)->comment('角色表主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('operator_id')->default(0)->comment('操作人（管理员）id')->index();
                $table->string('role_name')->default('')->comment('角色名称');
                $table->tinyInteger('type')->default(1)->comment('类型 1-功能角色 2-数据角色（暂无）');
                $table->text('auth_rule')->default('')->comment('权限规则json数据');
                $table->tinyInteger('identity_type')->default(1)->comment('身份类型 1-常规 2-店长 3-店员');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("alter table `enterprise_role` comment '企业角色表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_role');
    }
}
