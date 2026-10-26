<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddDeleteTypeToMemberUserRelation extends Migration
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
                $table->unsignedInteger('delete_type')->default(0)->after('union_id')->comment('删除方式  1 删除外部联系人 2 删除跟进成员');
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
                $table->dropColumn('delete_type');
            });
        }
    }
}
