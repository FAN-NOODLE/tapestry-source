<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateMemberUserRelationVisitTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('member_user_relation_visit')) {
            Schema::create('member_user_relation_visit', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->char('mem_uid', 32)->default('')->comment('UC人员（导购）id')->index();
                $table->unsignedInteger('count')->default(0)->comment('员工标记外部联系人访问人数');
                $table->unsignedInteger('result_time')->default(0)->comment('员工标记外部联系人访问人数查询时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `member_user_relation_visit` comment '员工标记外部联系人访问人数表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_user_relation_visit');
    }
}
