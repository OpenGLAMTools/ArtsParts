const Artwork = Vue.component('Artwork', {
  template: `
  <div class="modal fade" tabindex="-1" role="dialog" aria-labelledby="artworkeditLabel">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
      Edit
      </div>
      <div class="modal-body">
      Hallo
      <img :src=imagepath />
      </div>
    </div>
  </div>`,
  computed:{
    imagepath: function(){
      return "/img/"+this.iid+"/"+this.cid+"/"+this.aid;
    }
  },
  props: ['collection', 'iid', 'cid', 'aid']
})
const Collection = Vue.component('Collection', {
  template: `<div><h3>{{institution.collections[cid].name}}</h3>
  <ul>
  <li v-for="a in institution.collections[cid].artworks">
  <router-link 
        data-toggle="modal" data-target="#artworkedit"
        :to="{name: 'artwork', params:{iid:iid, cid:cid, aid:a.id}}">{{a.name}}</router-link>
  </li>
  </ul>
  <router-view :collection=collection id="artworkedit"></router-view>
  </div>`,
  computed: {
    collection: function () {
      return this.institution.collections[this.coll];
    },
    artworks: function () {
      return makeID(this.institution.collections[this.coll].artworks);
    }
  },
  props: ['institution', 'iid', 'cid']
})
const Institution = Vue.component('Institution', {
  template: `<div class=""><div class="row">
  <div class="col-md-2">
    <ul class="nav nav-pills nav-stacked">
    <router-link 
      v-for="c in data[iid].collections" :key="c.id"
      tag="li"
      class="presentation"
      active-class="active" 
      :to="{name: 'collection', params:{iid:iid, cid:c.id}}"><a>{{c.name}}</a>
    </router-link>
    </ul>
  </div>
    <div class="col-md-10">
    <router-view :institution=data[iid] :iid=iid></router-view>
    </div>
  </div>
  </div>`,
  props: {
    iid: {
      type: String,
      default: "dflt"
    },
    data: {
      type: Object,
      default: function () {
        return {
          dflt: {
            id: "dflt",
            name: "",
            collections: [
              {
                id: "",
                name: "",
              }
            ]
          }
        }
      }
    }
  },
  data: function () {
    return {
    }
  }
})
const Institutions = Vue.component('Institutions', {
  template: `<div class="container">
  
<ul class="nav nav-tabs">
<router-link 
  v-for="i in data" :key="i.id"
  tag="li"
  class="presentation"
  active-class="active" 
  :to="{name: 'institution', params:{iid: i.id}}"><a>{{i.name}}</a>
</router-link>
</ul>
<router-view :data=data></router-view>
  </div>`,
  data: function () {
    return {
    };
  },
  props: {
    data: {
      type: Object,
      default: [
        {
          id: "",
          name: "",
          collections: [
            {
              id: "",
              name: "",
            }
          ]
        }
      ]
    }
  }
});

const routes = [
  {
    name: 'institutions',
    path: '/',
    component: Institutions,
    children: [
      {
        name: 'institution',
        path: ':iid',
        component: Institution,
        props: true,
        children: [
          {
            name: 'collection',
            path: ':cid',
            props: true,
            component: Collection,
            children: [
              {
                name: 'artwork',
                path: ':aid',
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
    data: {
      dflt: {
        id: "dflt",
        name: "",
        collections: [
          {
            id: "",
            name: "",
            dflt: {
              id: "",
              name: "",
              dflt: {
                id: "",
                name: ""
              }
            }
          }
        ]
      }
    },
    appdata: [{ id: "", name: "" }]
  },
  methods: {
    fetchData: function () {
      this.dataLoaded = false;
      this.$http.get('/data/admin').then(
        (res) => {
          var data = {};
          res.body.forEach(function (el, ind, arr) {
            data[el.id] = el;
          });
          this.data = data;
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

function makeID(arr) {
  var out = {};
  arr.forEach(function (el, ind, arr) {
    out[el.id] = el;
  })
  return out;
}