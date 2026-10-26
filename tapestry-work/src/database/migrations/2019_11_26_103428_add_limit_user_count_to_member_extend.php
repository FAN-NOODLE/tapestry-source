<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddLimitUserCountToMemberExtend extends Migration
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
                $table->unsignedInteger('limit_user_count')->default(0)->after('allocation_status')->comment('分配客户上限');
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
                $table->dropColumn('limit_user_count');
            });
        }
    }
}
