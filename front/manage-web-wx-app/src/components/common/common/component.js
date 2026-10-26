import { basic } from '../mixins/basic'

function mapKeys(source, target, map = {}) {
  for (const key in map) {
    if (Object.hasOwnProperty.call(map, key)) {
      target[map[key]] = source[key]
    }
  }
}
export function MpComponent(mpOptions) {
  const options = {}
  mapKeys(mpOptions, options, {
    externalClasses: 'externalClasses',
    data: 'data',
    properties: 'properties',
    behaviors: 'behaviors',
    observers: 'observers',
    methods: 'methods',
    lifetimes: 'lifetimes',
    pageLifetimes: 'pageLifetimes'
  })
  // add default externalClasses
  options.externalClasses = options.externalClasses || []
  options.externalClasses.push('custom-class')
  // add default behaviors
  options.behaviors = options.behaviors || []
  options.behaviors.push(basic)
  // add relations
  if (mpOptions.relation) {
    options.relations = mpOptions.relation.relations
    options.behaviors.push(mpOptions.relation.mixin)
  }

  // map field to form-field behavior
  if (mpOptions.field) {
    options.behaviors.push('wx://form-field')
  }

  options.options = {
    multipleSlots: true,
    addGlobalClass: true
  }
  Component(options)
}
