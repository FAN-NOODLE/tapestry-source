// import { mount, createLocalVue } from "@vue/test-utils"
// import Vuex from "vuex"
// import { Button } from "element-ui"

// const localVue = createLocalVue()
// localVue.component(Button.name, Button)
// localVue.use(Vuex)
// const store = new Vuex.Store({});
// // const wrapper = mount(Demo, { store, localVue })

// describe("storespec.vue", () => {


//   it("默认状态值为空", () => {
//     expect(wrapper.vm.$data.userinfo).toBe(undefined)
//   })
  
//   it('点击按钮', () => {
//     wrapper.find(Button).vm.$emit("click")
//     wrapper.setData({
//       userinfo: {
//         username: '齐文豪',
//         age: 20,
//         high: '165'
//       }
//     })
//     jest.useFakeTimers()
//     console.log('wrapper.find(".inner")', wrapper.find(".inner"))
//     expect(wrapper.find(".inner").text()).toMatch("齐文豪")
//   })

//   // it("由子组件触发click事件，更新count为22", () => {
//   //   wrapper.setData({ count: 21 })
//   //   wrapper.find(Button).vm.$emit("click")
//   //   expect(wrapper.vm.$data.count).toBe(22)
//   // })

//   // it("由子组件内的button触发点击事件，更新count为11", () => {
//   //   wrapper.setData({ count: 10 })
//   //   const elBtn = wrapper.find(Button)
//   //   elBtn.find("button").trigger("click")
//   //   expect(wrapper.vm.$data.count).toBe(11)
//   // })
// })
