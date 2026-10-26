<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSourceTypeToMemberUserRelation extends Migration
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
                $table->unsignedInteger('source_type')->default(1)->after('delete_type')->comment('添加方式 1 导购主动添加  2 支付即导购添加  3 社群添加');
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
                $table->dropColumn('source_type');
            });
        }
    }
}
