<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddReleaseTagToEnterpriseThirdPartAccountTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->string('release_tag')->default('')->after('status')->comment('小程序发布版本分支tag');
                $table->string('release_branch')->default('')->after('status')->comment('小程序发布版本分支');
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
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->dropColumn('release_tag');
            });
        }
    }
}
