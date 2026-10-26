<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddStoreStaffIdInCrsOrderTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('crs_order')) {
            Schema::table('crs_order', function (Blueprint $table) {
                if (!Schema::hasColumn('crs_order', 'store_staff_id')) {
                    $table->unsignedInteger('store_staff_id')->default(0)->comment('分享人ID')->after('user_id')->index();
                }
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
        //
    }
}
