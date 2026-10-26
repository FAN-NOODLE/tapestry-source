<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseAdminRoleTable extends Migration
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
                $table->increments('id')->comment('id-自增');
                $table->unsignedInteger('enterprise_admin_id')->default(0)->comment('管理员id');
                $table->unsignedInteger('role_id')->default(0)->comment('角色id');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `enterprise_admin_role` COMMENT '管理员角色关系表'");
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
