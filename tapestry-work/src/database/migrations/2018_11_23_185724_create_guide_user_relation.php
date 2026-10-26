<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateGuideUserRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('guide_user_relation')) {
            Schema::create('guide_user_relation', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->char('mem_uid', 32)->default('')->comment('UC人员（导购）id')->index();
                $table->string('guide_id', 80)->default('')->index()->comment('导购统一id');
                $table->string('union_id', 80)->default('')->index()->comment('用户统一id');
                $table->unsignedInteger('user_id')->default(0)->comment('用户id')->index();
                $table->unsignedInteger('relation_at')->default(0)->comment('绑定时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `guide_user_relation` comment '客户导购关联关系表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('guide_user_relation');
    }
}
