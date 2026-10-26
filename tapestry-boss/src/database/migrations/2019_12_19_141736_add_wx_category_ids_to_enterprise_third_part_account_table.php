<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddWxCategoryIdsToEnterpriseThirdPartAccountTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_third_part_account', function (Blueprint $table) {
            $table->text('wx_category')->after('enterprise_type')->comment('小程序已设置的账号类目，json存储');
            $table->string('wx_category_ids', 191)
                ->default('')
                ->after('wx_category')
                ->comment('小程序账号类目id集合，json存储');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_third_part_account', function (Blueprint $table) {
            $table->dropColumn('wx_category');
            $table->dropColumn('wx_category_ids');
        });
    }
}
