<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\Resource;

class EnterpriseResource extends Resource
{
    /**
     * @param Request $request
     * @return array
     */
    public function toArray($request)
    {
        $identifier = 'contact';
        if ($this->contact_status == 3) {
            $identifier = 'crs-shop';
        }
        return [
            'id'              => $this->id,
            'hash'            => $this->hash,
            'name'            => $this->name,
            'logo'            => $this->logo,
            'contact'         => $this->contact,
            'mobile'          => $this->mobile,
            'address'         => $this->address,
            'crs_status'      => $this->crs_status,
            'store_status'    => $this->store_status,
            'erp_status'      => $this->erp_status,
            'crm_status'      => $this->crm_status,
            'contact_status'  => $this->contact_status,
            'identifier'      => $identifier, // cdp需要应用标识调uc
            'corp_id'         => $this->corp_id,
            'domain'          => $this->domain,
            'version'         => $this->version,
            'status'          => $this->status,
            'edit_at'         => $this->edit_at,
            'created_at'      => $this->created_at,
            'return_address'  => $this->return_address,
            'start_fee'       => $this->start_fee,
            'proportion'      => $this->proportion,
            'lock_time'       => $this->lock_time,
            'separate_switch' => $this->separate_switch,
        ];
    }
}
