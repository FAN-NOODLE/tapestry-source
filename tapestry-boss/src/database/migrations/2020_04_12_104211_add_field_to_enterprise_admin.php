<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldToEnterpriseAdmin extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_admin', function (Blueprint $table) {
            //
            $table->string('mobile')->default('')->after('salt')->comment('手机号码');
            $table->string('email')->default('')->after('salt')->comment('邮箱（目前作为超管账号');
            $table->string('global_roaming')->default('')->after('salt')->comment('国际区号');
            $table->tinyInteger('is_super')->default(2)->after('role')->comment('是否超管: 1-是 2-不是');
            $table->tinyInteger('range_type')->default(1)->comment('账号范围 1-不限 2-指定范围');
            $table->text('admin_available_system')->comment('账号可用系统,json字符存储[1-CRS后台 2-店长小程序 3-天擎CDP]');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_admin', function (Blueprint $table) {
            //
            $table->dropColumn('mobile');
            $table->dropColumn('email');
            $table->dropColumn('global_roaming');
            $table->dropColumn('is_super');
        });
    }
}
