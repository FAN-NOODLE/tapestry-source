# Code lookup map

Keep full source in this repository and business context in Notion. Use Notion feature IDs and repository-relative paths to retrieve specific files instead of traversing the entire repository for each question.

## Front-end projects

- `front/login-web-admin/`: login and platform administration entry.
- `front/web-admin/`: business administration UI.
- `front/manage-web-wx-app/`: employee mini-program.
- `front/tapestry-work-web-h5/`: enterprise-WeChat H5 UI.

## Services

- `micro-download/`
- `tapestry-authority/`
- `tapestry-boss/`
- `tapestry-gw-assistant/`
- `tapestry-gw-console/`
- `tapestry-gw-scrm/`
- `tapestry-gw-work-scrm/`
- `tapestry-micro-scrm/`
- `tapestry-work/`

## Example: F16-007 order cancellation

1. H5 entry and confirmation: `front/tapestry-work-web-h5/src/pages/online-order/detail.vue`, `handleCancelOrder` / `handleCancelOrderConfirm`.
2. H5 request: `front/tapestry-work-web-h5/src/api/online-order/index.ts`, `cancelOrderApi`.
3. Gateway: `tapestry-gw-work-scrm/internal/app/http/handle/online_warehouse/order.go`, `OrderHandle.Cancel`.
4. Service handler: `tapestry-micro-scrm/internal/app/http/handle/online_warehouse/order.go`, `Handle.Cancel`.
5. Business checks and updates: `tapestry-micro-scrm/internal/app/service/online_warehouse/order.go`, `OrderSvc.Cancel`, original lines 1542–1567.
6. Status constants: `tapestry-micro-scrm/internal/app/models/micro_scrm/online_warehouse/order.go`.

Verify ownership, order status and `is_export` separately. Normal export and shipping-import-template generation have different effects. Check the Notion feature page for evidence and limitations.
