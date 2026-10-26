export function useChildren(name, onEffect) {
  const path = '../'.concat(name, '/index')

  return {
    relations: {
      [path]: {
        type: 'descendant',
        linked: function (target) {
          onEffect && onEffect.call(this, target)
        },
        linkChanged: function (target) {
          onEffect && onEffect.call(this, target)
        },
        unlinked: function (target) {
          onEffect && onEffect.call(this, target)
        }
      }
    },
    mixin: Behavior({
      lifetimes: {
        created() {
          const _this = this
          Object.defineProperty(this, 'children', {
            get: function () {
              return _this.getRelationNodes(path) || []
            }
          })
        }
      }
    })
  }
}

export function useParent(name, onEffect) {
  const path = '../'.concat(name, '/index')
  return {
    relations: {
      [path]: {
        type: 'ancestor',
        linked: function () {
          onEffect && onEffect.call(this)
        },
        linkChanged: function () {
          onEffect && onEffect.call(this)
        },
        unlinked: function () {
          onEffect && onEffect.call(this)
        }
      }
    },
    mixin: Behavior({
      lifetimes: {
        created() {
          const _this = this
          Object.defineProperty(this, 'parent', {
            get: function () {
              return _this.getRelationNodes(path)[0]
            }
          })
          // console.log();
          Object.defineProperty(this, 'index', {
            // @ts-ignore
            get: function () {
              return _this.parent.children.indexOf(_this)
            }
          })
        }
      }
    })
  }
}
