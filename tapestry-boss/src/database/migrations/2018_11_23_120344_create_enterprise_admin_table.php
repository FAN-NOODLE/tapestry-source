<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseAdminTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_admin')) {
            Schema::create('enterprise_admin', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->string('username',40)->default('')->comment('用户名');
                $table->string('account',20)->default('')->comment('账户');
                $table->string('password',40)->default('')->comment('密码');
                $table->string('salt',6)->default('')->comment('盐');
                $table->unsignedBigInteger('role')->default(1)->comment('角色（1：系统管理员）');
                $table->unsignedInteger('last_login_time')->default(0)->comment('最后一次登录时间');
                $table->unsignedTinyInteger('status')->default(1)->comment('状态（1:正常，2:禁用）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id_index');
                $table->unique('account', 'account_unique');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `enterprise_admin` comment '企业管理员表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_admin');
    }
}
