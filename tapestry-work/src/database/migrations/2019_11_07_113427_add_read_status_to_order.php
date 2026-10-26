<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddReadStatusToOrder extends Migration
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
                $table->unsignedInteger('read_status')->default(1)->after('other_fee')->comment('已读状态  1 未读 2 已读');
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
        if (Schema::hasTable('order')) {
            Schema::table('order', function (Blueprint $table) {
                $table->dropColumn('read_status');
            });
        }
    }
}
