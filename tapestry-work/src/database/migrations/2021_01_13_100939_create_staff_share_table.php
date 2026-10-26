<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStaffShareTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('staff_share')) {
            Schema::create('staff_share', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment("主键ID");
                $table->unsignedBigInteger("enterprise_id")->default(0)->comment("企业ID");
                $table->unsignedBigInteger("store_id")->default(0)->comment("门店ID");
                $table->unsignedTinyInteger("type")->default(1)->comment("1:内购 2:导购 10:合伙人");
                $table->unsignedBigInteger("user_id")->default(0)->comment("用户ID");
                $table->unsignedBigInteger("parent_user_id")->default(0)->comment("父级用户ID");
                $table->string("mem_uid", 32)->default('')->comment("员工ID");
                $table->unsignedBigInteger("worker_id")->default(0)->comment("员工workerid");
                $table->unsignedInteger('level')->default(0)->comment('计数');
                $table->string('qr_code')->default('')->comment('二维码路径');
                $table->unsignedTinyInteger('binding_type')->default(0)->comment('1:绑定 2:减除绑定 3:客户分享,分享人无绑定');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index(['enterprise_id'], 'enterprise_id');
                $table->index(['mem_uid'], 'mem_uid');
                $table->index(['user_id'], 'user_id');
            });
        }

        DB::statement("alter table `staff_share` comment '员工内购分享链路表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('staff_share');
    }
}
