import { transition } from '../mixins/transition'

Component({
  externalClasses: [
    'custom-class',
    'enter-class',
    'enter-active-class',
    'enter-to-class',
    'leave-class',
    'leave-active-class',
    'leave-to-class',
    'close-icon-class'
  ],
  options: {
    virtualHost: false,
    multipleSlots: true,
    styleIsolation: 'apply-shared'
  },
  behaviors: [transition(false)],
  properties: {
    zIndex: {
      type: Number,
      value: 521
    },
    position: {
      type: String,
      value: 'center',
      observer: 'observeClass'
    },
    overlay: {
      type: Boolean,
      value: true
    },
    duration: {
      type: Number,
      value: 300
    },
    round: {
      type: Boolean,
      value: false
    },
    customStyle: {
      type: String,
      value: ''
    },
    overlayStyle: {
      type: String,
      value: ''
    },
    closeOnClickOverlay: {
      type: Boolean,
      value: true
    },
    closeable: {
      type: Boolean,
      value: false
    },
    closeIcon: {
      type: String,
      value: 'cross'
    },
    closeIconPosition: {
      type: String,
      value: 'top-right'
    },
    safeAreaInsetBottom: {
      type: Boolean,
      value: false
    },
    safeAreaInsetTop: {
      type: Boolean,
      value: false
    },
    safeAreaTabBar: {
      type: Boolean,
      value: false
    },
    lockScroll: {
      type: Boolean,
      value: true
    },
    rootPortal: {
      type: Boolean,
      value: false
    },
    transition: {
      type: String,
      observer: 'observeClass'
    }
  },
  data: {},
  lifetimes: {
    created() {
      this.observeClass()
    }
  },
  methods: {
    onClickCloseIcon() {
      this.triggerEvent('close')
    },
    onClickOverlay() {
      this.triggerEvent('click-overlay')
      if (this.data.closeOnClickOverlay) {
        this.triggerEvent('close')
      }
    },
    observeClass() {
      const { transition, position, duration } = this.data
      var updateData = {
        name: transition || position
      }
      if (transition === 'none') {
        updateData.duration = 0
        this.originDuration = duration
      } else if (
        this.originDuration !== null &&
        this.originDuration !== undefined
      ) {
        updateData.duration = this.originDuration
      }
      this.setData(updateData)
    },
    noop() {}
  }
})
