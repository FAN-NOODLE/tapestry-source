const { getCurrentPage } = require('@/components/common/common/utils')


export function toggleTab() {
  const page = getCurrentPage()
  page.toggleTab()
}
