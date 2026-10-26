<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use \Illuminate\Support\Facades\DB;

class CreateOrderStatusTraceTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('order_status_trace')) {
            Schema::create('order_status_trace', function (Blueprint $table) {
                $table->increments('id');
                $table->string('order_no', 30)->default('')->comment('订单号')->index();
                $table->string('event', 200)->default('')->comment('订单变化事件');
                $table->string('event_desc', 200)->default('')->comment('订单变化事件描述');
                $table->unsignedInteger('from')->default(0)->comment('变化之前状态');
                $table->unsignedInteger('to')->default(0)->comment('变化之后状态');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('change_at')->default(0)->comment('状态变化时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `order_status_trace` comment '订单状态变化表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('order');
    }
}
