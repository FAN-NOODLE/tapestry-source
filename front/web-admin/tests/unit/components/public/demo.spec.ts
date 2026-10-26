import { mount, createLocalVue } from "@vue/test-utils"
import Demo from "@/components/public/demo/index.vue"
import { Button } from "element-ui"

const localVue = createLocalVue()
localVue.component(Button.name, Button)
const wrapper = mount(Demo, { localVue })

describe("demo.vue", () => {
  it("默认状态值为0", () => {
    expect(wrapper.vm.$data.count).toBe(0)
  })

  // it("修改count值为20", () => {
  //   wrapper.setData({ count: '20' })
  //   console.log('wrapper.find("span").text()', wrapper.find("span").text())
  //   expect(wrapper.find("span").text()).toMatch("20")
  // })

  it("由子组件触发click事件，更新count为22", () => {
    wrapper.setData({ count: 21 })
    wrapper.find(Button).vm.$emit("click")
    expect(wrapper.vm.$data.count).toBe(22)
  })

  it("由子组件内的button触发点击事件，更新count为11", () => {
    wrapper.setData({ count: 10 })
    const elBtn = wrapper.find(Button)
    elBtn.find("button").trigger("click")
    expect(wrapper.vm.$data.count).toBe(11)
  })
})
