<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddStoreIdToTaskMemberRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('task_member_relation')) {
            Schema::table('task_member_relation', function (Blueprint $table) {
                $table->string('store_id')->default('')->after('task_id')->comment('门店id');
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
        if (Schema::hasTable('task_member_relation')) {
            Schema::table('task_member_relation', function (Blueprint $table) {
                $table->dropColumn('store_id');
            });

        }
    }
}
