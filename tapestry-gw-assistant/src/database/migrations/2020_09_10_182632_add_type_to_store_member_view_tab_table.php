<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddTypeToStoreMemberViewTabTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('type', 'store_member_view_tab')) {
            Schema::table('store_member_view_tab', function (Blueprint $table) {
                $table->unsignedInteger('type')->default(1)->after('store_worker_id')->comment('类型 1:查推广订单tab 2 待接订单 3:顾客列表');
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
        Schema::table('store_member_view_tab', function (Blueprint $table) {
            $table->dropColumn(['type']);
        });
    }
}
