<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddObjectExtraIdToStaffRightTable extends Migration
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
                $table->text('object_extra_id')->default('')->after('obj_id')->comment('部门/员工当前有权限的员工id集合');
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
                $table->dropColumn('object_extra_id');
            });
        }
    }
}
