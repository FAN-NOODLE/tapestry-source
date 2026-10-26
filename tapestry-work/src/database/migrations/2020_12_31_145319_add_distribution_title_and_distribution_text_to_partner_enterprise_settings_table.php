<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddDistributionTitleAndDistributionTextToPartnerEnterpriseSettingsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('partner_enterprise_settings')) {
                if (!Schema::hasColumn('partner_enterprise_settings', 'distribution_title')) {
                    Schema::table('partner_enterprise_settings', function (Blueprint $table) {
                        $table->string('distribution_title', 40)->default('')->comment('二级分销公告标题');
                        $table->text('distribution_text')->default('')->comment('二级分销公告内容');
                    });
                }
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('partner_enterprise_settings', function (Blueprint $table) {
            $table->dropColumn('distribution_title');
            $table->dropColumn('distribution_text');
        });
    }
}
