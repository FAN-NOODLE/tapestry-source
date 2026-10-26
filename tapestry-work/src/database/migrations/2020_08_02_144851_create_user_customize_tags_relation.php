<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateUserCustomizeTagsRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('user_customize_tags_relation', function (Blueprint $table) {
            $table->increments('id')->comment('自增主键');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
            $table->unsignedBigInteger('worker_id')->default(0)->comment('员工(store_worker)ID')->index();
            $table->string('union_id', 50)->default('')->comment('uc外部联系人union_id')->index();
            $table->unsignedBigInteger('user_tags_id')->default(0)->comment('自定义标签(user_customize_tags)ID')->index();
            $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间')->index();
        });
        // 表注释
        DB::statement("ALTER TABLE `user_customize_tags_relation` comment '用户自定义标签关联表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('user_customize_tags_relation');
    }
}
