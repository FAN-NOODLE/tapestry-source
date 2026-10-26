<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateUserStoreWorkerRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('user_store_worker_relation')) {
            Schema::create('user_store_worker_relation', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment("主键");
                $table->unsignedBigInteger("user_id")->default(0)->comment("用户ID");
                $table->unsignedBigInteger("enterprise_id")->default(0)->comment("企业ID");
                $table->unsignedBigInteger("store_id")->default(0)->comment("门店ID");
                $table->string("mem_uid")->default('')->comment("UC员工ID");
                $table->unsignedBigInteger("worker_id")->default(0)->comment("员工ID");
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index(['user_id', 'enterprise_id', 'store_id', 'worker_id'], 'user_store_worker_index');
            });
        }
        DB::statement("alter table `update_worker_record` comment '用户门店员工绑定关联表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('user_store_worker_relation');
    }
}
