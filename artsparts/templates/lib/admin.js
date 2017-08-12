const Artwork = Vue.component('Artwork', {
  template: `
  <div class="ui modal">
  <i class="close icon"></i>
  <div class="header">
    {{artwork.name}}
  </div>
  <div class="image content">
    <div class="ui medium image">
    
      <img :src=imagepathhuge />
    
    </div>
    <div class="description">
    <p>Edit the data to that artwork. Be sure to save your changes.</p  >
      <div class="ui form">
          <div class="field">
            <label for="Name">Title</label>
            <input type="text" class="form-control" id="Name" v-model="artwork.name">
          </div>
          <div class="field">
            <label for="Timestamp">Timestamp</label>
            <input  class="form-control" id="Timestamp" placeholder="YYYYMMDDHHMM" v-model="artwork.timestamp">
          </div>
          <div class="field">
            <label for="HashTag">HashTag</label>
            <input  class="form-control" id="HashTag"  v-model="artwork.hashtag">
          </div>
          <div class="field" v-for="d in meta">
            <label :for=d>{{d}}</label>
            <input  class="form-control" :id="d"  v-model="artwork.meta[d]">
          </div>
          <div class="field">
            <label for="Desc">Description</label>
            <textarea id="Desc" class="form-control" rows="5" v-model="artwork.description"></textarea>
          </div>
          
        </div>     
    </div>
  </div>
  <div class="actions">
    <div class="ui black deny button">
      Close and Reload
    </div>
    <div class="ui positive right labeled icon button" v-on:click="save">
      Save changes
      <i class="checkmark icon"></i>
    </div>
  </div>
</div>

  `,
  data: function(){
    return {
      artwork: this.collection.artworks[this.aindex],
      meta : ['artist','date','link']
    }
  },
  computed:{
    imagepathmedium: function(){
      return '/img/'+this.iid+'/'+this.cid+'/'+this.artwork.id+"?size=medium";
    },
    imagepathhuge: function(){
      return '/img/'+this.iid+'/'+this.cid+'/'+this.artwork.id+"?size=massive";
    }
  },
  methods: {
    update: function(){
      this.artwork =  this.collection.artworks[this.aindex]
    },
    save: function(){
      this.$http.post('/data/'+this.iid+'/'+this.cid+'/'+this.artwork.id, this.artwork).then(response => {
        $('#artworkedit').modal('hide')
    // success callback
    // console.log("Artwork is safed");
    }, response => {
      // error callback
      // console.log("There was an error");
    });
    }
  },
  watch:{
    '$route': 'update'
  }, 
  props: ['collection', 'iid', 'cid', 'aindex']
})
const Collection = Vue.component('Collection', {
  template: `<div>
  <table class="ui celled striped very compact table">
 <thead>
  <tr><th>ID</th><th>Name</th><th>Description</th><th>Timestamp</th></tr>
</thead>
  <tbody>
  <tr v-for="(a,i) in institution.collections[cid].artworks">
  <td v-on:click="modal"><router-link 
        :to="{name: 'artwork', params:{iid:iid, cid:cid, aindex:i}}">{{a.id}}<br>
        <img :src=imgPath(a.id)>
        </router-link>
        </td>
<td v-on:click="modal"><router-link 
        
        :to="{name: 'artwork', params:{iid:iid, cid:cid, aindex:i}}">{{a.name}}</router-link></td>
<td>{{a.description}}</td>
<td>{{renderTS(a.timestamp)}}
</td>                        
  </tr>
  </tbody>
  </table>
  
  <router-view :collection=institution.collections[cid] id="artworkedit"></router-view>
  </div>`,
  computed: {
    collection: function () {
      return this.institution.collections[this.coll];
    },
    artworks: function () {
      return makeID(this.institution.collections[this.coll].artworks);
    }
  },
  methods: {
    modal: function(){
       $('.ui.modal').modal('show');
    },
    imgPath: function(aid){
      return "/img/"+this.iid+"/"+this.cid+"/"+aid+"?size=small";
    },
    renderTS: function(ts){
      return ts.substr(6,2)+"."+ts.substr(4,2)+"."+ts.substr(0,4)+" "+ts.substr(8,2)+":"+ts.substr(10,2)
    }
  },
  props: ['institution', 'iid', 'cid']
})
const Institution = Vue.component('Institution', {
  template: `<div class="ui stacked segment grid">
  <div class="ui vertical pointing menu four wide column">
    <router-link 
      v-for="c in data[iid].collections" :key="c.id"
      class="item"
      active-class="active" 
      :to="{name: 'collection', params:{iid:iid, cid:c.id}}">{{c.name}}
    </router-link>
  </div>
    <div class="twelve wide column">
    <router-view :institution=data[iid] :iid=iid></router-view>
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
  template: `<div class="ui container">
  <div class="ui pointing menu">
<router-link 
  v-for="i in data" :key="i.id"
  class="item"
  active-class="active" 
  :to="{name: 'institution', params:{iid: i.id}}">{{i.name}}
</router-link>
</div>
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
                path: ':aindex',
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
  template: `<div class="ui segment">
  <ul>
  <li><router-link to="/">
  Administrate your Institutions</router-link></li>
  </ul>
  <router-view :data=data></router-view></div>`,
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