<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCustomerGroupRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('customer_group_relation')) {
            Schema::create('customer_group_relation', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('customer_group_id')->default(0)->comment('用户分组id')->index();
                $table->unsignedBigInteger('user_id')->default(0)->comment('用户id');
                $table->string('ec_userid')->default('')->comment('uc用户id');
                $table->text('ec_info')->default('')->comment('uc信息');
                $table->string('receipt_handle')->default('')->comment('队列标识');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
        }

        // 表注释
        \Illuminate\Support\Facades\DB::statement("ALTER TABLE `customer_group_relation` comment '用户分组和用户关系表'");

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('customer_group_relation');
    }
}
