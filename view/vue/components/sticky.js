var Sticky = Vue.extend(
    {
        template: '<div :style="{height:height+\'px\',zIndex:zIndex}"> <div :class="className" :style="{top:stickyTop+\'px\',zIndex:zIndex,position:position,width:width,height:height+\'px\',backgroundColor:\'white\'}"><slot><div>sticky</div></slot></div></div>',
        name: 'Sticky',
        props: {
            stickyTop: {
                type: Number,
                default: 0
            },
            height: {
                type: Number,
                default: 50
            },
            zIndex: {
                type: Number,
                default: 1
            },
            className: {
                type: String
            }
        },
        data() {
            return {
                active: false,
                position: '',
                currentTop: '',
                width: undefined,
                child: null,
                stickyHeight: 0
            }
        },
        methods: {
            sticky() {
                if (this.active) {
                    return;
                }
                this.position = 'fixed';
                this.active = true;
                this.width = this.width + 'px';
            },
            reset() {
                if (!this.active) {
                    return
                }
                this.position = '';
                this.width = 'auto';
                this.active = false
            },
            handleScroll() {
                this.width = this.$el.getBoundingClientRect().width;
                const offsetTop = this.$el.getBoundingClientRect().top;
                if (offsetTop <= this.stickyTop) {
                    this.sticky();
                    return;
                }
                this.reset();
            }
        },
        mounted() {
            window.addEventListener('scroll', this.handleScroll);
        },
        destroyed() {
            window.removeEventListener('scroll', this.handleScroll);
        }
    }
);
