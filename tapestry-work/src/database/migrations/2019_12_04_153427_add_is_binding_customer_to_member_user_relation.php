<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIsBindingCustomerToMemberUserRelation extends Migration
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
                $table->unsignedInteger('is_binding_customer')->default(1)->after('source_type')->comment('是绑定当前导购的客户 1 是 2否');
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
                $table->dropColumn('is_binding_customer');
            });
        }
    }
}
