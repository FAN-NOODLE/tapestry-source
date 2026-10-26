<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIsMyOrderToStaffOrderRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('staff_order_relation')) {
            Schema::table('staff_order_relation', function (Blueprint $table) {
                $table->unsignedBigInteger('is_my_order')->default(0)->after('order_money')->comment('是否是员工自己的订单（1=是、2=否）');
                $table->string('dp_id',32)->default('')->after('order_money')->comment('部门id');

            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('staff_order_relation')) {
            Schema::table('staff_order_relation', function (Blueprint $table) {
                $table->dropColumn('is_my_order');
                $table->dropColumn('dp_id');
            });
        }
    }
}
