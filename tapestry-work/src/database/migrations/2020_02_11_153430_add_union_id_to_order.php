<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddUnionIdToOrder extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('order')) {
            Schema::table('order', function (Blueprint $table) {
                $table->string('union_id')->default('')->after('user_id')->comment('用户union_id');
                $table->string('open_id')->default('')->after('user_id')->comment('用户open_id');
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
        if (Schema::hasTable('user_coupon')) {
            Schema::table('user_coupon', function (Blueprint $table) {
                $table->dropColumn('union_id');
                $table->dropColumn('open_id');
            });
        }
    }
}
