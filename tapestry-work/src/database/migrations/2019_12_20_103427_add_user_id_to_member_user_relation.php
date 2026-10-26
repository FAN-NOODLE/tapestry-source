<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddUserIdToMemberUserRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('member_user_relation')) {
            Schema::table('member_user_relation', function (Blueprint $table) {
                $table->unsignedBigInteger('user_id')->default(0)->after('source_type')->comment('用户商城id');
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
        if (Schema::hasTable('member_user_relation')) {
            Schema::table('member_user_relation', function (Blueprint $table) {
                $table->dropColumn('user_id');
            });
        }
    }
}
