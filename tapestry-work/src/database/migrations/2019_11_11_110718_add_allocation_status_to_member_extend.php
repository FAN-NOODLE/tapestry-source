<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAllocationStatusToMemberExtend extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->unsignedTinyInteger('allocation_status')->default(2)->after('avatar')->comment('分配状态：1、开启，2、关闭，3、强制开启，4、永久关闭');
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
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->dropColumn('allocation_status');
            });
        }
    }
}
