<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTagCustomerGroupRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('tag_customer_group_relation')) {
            Schema::create('tag_customer_group_relation', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('tag_group_id')->default('')->comment('标签组id');
                $table->string('tag_id')->default('')->comment('标签id');
                $table->string('custom_group_ids')->default('')->comment('用户分组ids，数组');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
        }

        // 表注释
        \Illuminate\Support\Facades\DB::statement("ALTER TABLE `tag_customer_group_relation` comment '用户分组和标签关系表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('tag_customer_group_relation');
    }
}
