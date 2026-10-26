<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddObjParentIdToStaffRightTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('staff_right')) {
            Schema::table('staff_right', function (Blueprint $table) {
                $table->string('obj_parent_id')->default(0)->after('obj_id')->comment('员工所在部门ID');
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
        if (Schema::hasTable('staff_right')) {
            Schema::table('staff_right', function (Blueprint $table) {
                $table->dropColumn('obj_parent_id');
            });
        }
    }
}
