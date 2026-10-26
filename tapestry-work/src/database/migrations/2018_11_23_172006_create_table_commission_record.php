<?php

use \Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableCommissionRecord extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('commission_record')) {
            Schema::create('commission_record', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->char('mem_uid', 32)->default('')->index()->comment('导购id');
                $table->unsignedInteger('amount_price')->default(0)->comment('订单总额，单位：分');
                $table->unsignedInteger('order_id')->default(0)->comment('订单id');
                $table->unsignedInteger('user_id')->default(0)->index()->comment('会员id');
                $table->unsignedInteger('charge_price')->default(0)->comment('佣金金额，单位：分');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `commission_record` COMMENT '佣金记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('commission_record');
    }
}
