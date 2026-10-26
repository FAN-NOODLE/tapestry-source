<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAchievementDescToStaffSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('staff_setting')) {
            Schema::table('staff_setting', function (Blueprint $table) {
                $table->text('achievement_desc')->default('')->after('settle_type')->comment('业绩说明');
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
        if (Schema::hasTable('staff_setting')) {
            Schema::table('staff_setting', function (Blueprint $table) {
                $table->dropColumn('achievement_desc');
            });
        }
    }
}
