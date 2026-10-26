<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddRoleToMemberExtend extends Migration
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
                $table->char('role')->default('')->after('mem_uid')->comment('UC员工角色');
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
                $table->dropColumn('role');
            });

        }
    }
}
