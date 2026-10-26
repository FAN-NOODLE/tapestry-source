Component({
  externalClasses: [
    'custom-class'
  ],
  properties: {
    show: Boolean,
    mask: Boolean,
    title: String,
    message: String,
    forbidClick: Boolean,
    zIndex: {
      type: Number,
      value: 1000
    },
    type: {
      type: String,
      value: 'text'
    },
    position: {
      type: String,
      value: 'middle'
    }
  }
})
