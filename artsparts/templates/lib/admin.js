const Foo = { 
  template: '<div>foosdfsdf</div>'
 }
const Bar = { template: '<div>bar</div>' }

const routes = [
  { path: '/foo', component: Foo },
  { path: '/bar', component: Bar }
]

const router = new VueRouter({
  routes 
})

const app = new Vue({
  router,
  template: `<div><router-link to="/foo">Go to Foo</router-link>
    <router-link to="/bar">Go to Bar</router-link><router-view></router-view></div>`,
  data: {
    dataLoaded: false,
    appdata: {}
  },
  methods: {
        fetchData: function () {
            this.dataLoaded = false;
            this.$http.get('/data/admin').then(
                (res) => {
                    this.appdata = res.body;                 
                    this.dataLoaded = true;
                }
            );
        }
  },
  created: function () {
        this.fetchData()
    }

}).$mount("#adminapp")