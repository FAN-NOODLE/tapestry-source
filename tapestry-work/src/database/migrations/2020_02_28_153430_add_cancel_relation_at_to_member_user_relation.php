<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCancelRelationAtToMemberUserRelation extends Migration
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
                $table->unsignedInteger('cancel_relation_at')->default(0)->after('delete_type')->comment('取消关系时间 员工删客户 或 客户删员工');
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
                $table->dropColumn('cancel_relation_at');
            });
        }
    }
}
