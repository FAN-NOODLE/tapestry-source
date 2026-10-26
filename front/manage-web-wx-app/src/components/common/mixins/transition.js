import { requestAnimationFrame } from '../common/utils'
import { isObj } from '../common/validator'

var getClassNames = function (name) {
  return {
    enter: 'mp-'
      .concat(name, '-enter mp-')
      .concat(name, '-enter-active enter-class enter-active-class'),
    'enter-to': 'mp-'
      .concat(name, '-enter-to mp-')
      .concat(name, '-enter-active enter-to-class enter-active-class'),
    leave: 'mp-'
      .concat(name, '-leave mp-')
      .concat(name, '-leave-active leave-class leave-active-class'),
    'leave-to': 'mp-'
      .concat(name, '-leave-to mp-')
      .concat(name, '-leave-active leave-to-class leave-active-class')
  }
}

export function transition(showDefaultValue) {
  return Behavior({
    properties: {
      show: {
        type: Boolean,
        value: showDefaultValue,
        observer: 'observeShow'
      },
      duration: {
        type: Object | Number,
        value: 300,
        observer: 'observeDuration'
      },
      name: {
        type: String,
        value: 'fade'
      }
    },
    data: {
      type: '',
      inited: false,
      display: false
    },
    lifetimes: {
      ready() {
        if (this.data.show === true) {
          this.observeShow(true, false)
        }
      }
    },
    methods: {
      observeShow(value, old) {
        if (value === old) return
        value ? this.enter() : this.leave()
      },
      enter() {
        const { name, duration } = this.data
        const classNames = getClassNames(name)
        const currentDuration = isObj(duration) ? duration.enter : duration
        if (this.status === 'enter') {
          return
        }
        this.status = 'enter'
        this.triggerEvent('before-enter')
        requestAnimationFrame(() => {
          if (this.status !== 'enter') return
          this.triggerEvent('enter')
          this.setData({
            inited: true,
            display: true,
            classes: classNames.enter,
            currentDuration
          })
          requestAnimationFrame(() => {
            if (this.status !== 'enter') return
            this.transitionEnded = false
            this.setData({
              classes: classNames['enter-to']
            })
          })
        })
      },
      leave() {
        const { name, duration } = this.data
        const classNames = getClassNames(name)
        const currentDuration = isObj(duration) ? duration.leave : duration

        if (this.status === 'leave') {
          return
        }
        this.status = 'leave'
        this.triggerEvent('before-leave')
        requestAnimationFrame(() => {
          if (this.status !== 'leave') return
          this.triggerEvent('leave')
          this.setData({
            classes: classNames.leave,
            currentDuration
          })
          requestAnimationFrame(() => {
            if (this.status !== 'leave') return
            this.transitionEnded = false
            setTimeout(() => {
              return this.onTransitionEnd()
            }, currentDuration)
            this.setData({
              classes: classNames['leave-to']
            })
          })
        })
      },
      onTransitionEnd() {
        if (this.transitionEnded) {
          return
        }
        this.transitionEnded = true
        this.triggerEvent('after-'.concat(this.status))
        const { show, display } = this.data
        if (!show && display) {
          this.setData({ display: false })
        }
      }
    }
  })
}
