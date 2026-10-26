<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateUpdateWorkerRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('update_worker_record')) {
            Schema::create('update_worker_record', function (Blueprint $table) {
                $table->bigInteger('id')->primary()->comment("主键");
                $table->bigInteger("enterprise_id")->default(0)->comment("企业ID");
                $table->string("domain")->default('')->comment("企业domain");
                $table->bigInteger("worker_id")->default(0)->comment("导购id");
                $table->text("result")->comment("错误原因");
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id');
            });
        }
        DB::statement("alter table `update_worker_record` comment '更新员工记录表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('update_worker_record');
    }
}
