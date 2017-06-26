const Artwork = Vue.component('Artwork', {
  template: `<div></div>`
})
const Collection = Vue.component('Collection', {
  template: `<div></div>`
})
const Institution = Vue.component('Institution', {
  template: `<div>
<ul>
<li v-for="coll in institution.collections"><router-link :to="{name: 'collection', params:{inst:inst.id, coll:coll.id}}">{{coll.name}}</router-link></li>
</ul>
<router-view :institution=institution></router-view>
  </div>`,
  props: ['data', 'inst'],
  data: function () {
    return {
      institution: {}
    }
  },
  methods: {
    getInstitution: function () {
      this.insitution = this.data[this.inst];
    }
  },
  created: function () {
    this.getInstitution();
  }
})
const Institutions = Vue.component('Institutions', {
  template: `<div>Your Institutions:
<ul>
<li v-for="i in data"><router-link :to="{name: 'institution', params:{inst: i.id}}">{{i.name}}</router-link></li>
</ul>
<router-view :data=data></router-view>
  </div>`,
  data: function () {
    return {
    };
  },
  props: ['data']
});

const routes = [
  {
    name: 'institutions',
    path: '/',
    component: Institutions,
    children: [
      {
        name: 'institution',
        path: '/:inst',
        component: Institution,
        props: true,
        children: [
          {
            name: 'collection',
            path: '/:coll',
            props: true,
            component: Collection,
            children: [
              {
                name: 'artwork',
                path: '/:artw',
                props: true,
                component: Artwork
              }
            ]
          }
        ]
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

const app = new Vue({
  router,
  template: `<div><router-link to="/">Institutions</router-link><router-view :data=data></router-view></div>`,
  data: {
    dataLoaded: false,
    appdata: [{ id: "", name: "" }]
  },
  computed: {
    data: function () {
      var inst = {};
      this.appdata.forEach(function (el, ind, arr) {
        inst[el.id] = el;
      });
      return inst;
    }
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