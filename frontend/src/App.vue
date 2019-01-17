<template>
  <v-app id="gopicamera">
    <div v-if="isLoggedIn">
      <v-navigation-drawer
        :clipped="$vuetify.breakpoint.mdAndUp"
        v-model="drawer"
        fixed
        app
      >
        <v-list>
          <template v-for="item in items">
            <v-layout
              v-if="item.heading"
              :key="item.heading"
              row
              align-center
            >
              <v-flex xs6>
                <v-subheader v-if="item.heading">
                  {{ item.heading }}
                </v-subheader>
              </v-flex>
            </v-layout>
            <v-list-group
              v-else-if="item.children && !item.requiresAdmin || item.requiresAdmin && (item.requiresAdmin == isAdmin)"
              v-model="item.model"
              :key="item.text"
              :prepend-icon="item.model ? item.icon : item['icon-alt']"
              append-icon=""
            >
              <v-list-tile slot="activator">
                <v-list-tile-content>
                  <v-list-tile-title>
                    {{ item.text }}
                  </v-list-tile-title>
                </v-list-tile-content>
              </v-list-tile>
              <v-list-tile v-for="(child, i) in item.children" :key="i" @click="$router.push(child.route)" :to="child.route" router>
                <v-list-tile-action v-if="child.icon">
                  <v-icon>{{ child.icon }}</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>
                  <v-list-tile-title>
                    {{ child.text }}
                  </v-list-tile-title>
                </v-list-tile-content>
              </v-list-tile>
            </v-list-group>
            <v-list-tile v-else-if="!item.requiresAdmin || item.requiresAdmin && (item.requiresAdmin == isAdmin)" :key="item.text" @click="$router.push(item.route)" :to="item.route" router>
              <v-list-tile-action>
                <v-icon>{{ item.icon }}</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>
                <v-list-tile-title>
                  {{ item.text }}
                </v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
          </template>
          
          <!--v-list-tile @click="menuLogout()">
            <v-list-tile-action>
              <v-icon>exit_to_app</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                Logout
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile -->
        </v-list>
      </v-navigation-drawer>
    </div>
    
    <v-toolbar
      :clipped-left="$vuetify.breakpoint.mdAndUp"
      color="blue darken-2"
      dark
      app
      fixed
    >
      <v-toolbar-title style="width: 300px" class="ml-0">
        <v-toolbar-side-icon @click.stop="drawer = !drawer" v-if="isLoggedIn"></v-toolbar-side-icon>
        <span><v-icon color="blue lighten-4">camera</v-icon> gopicamera</span>
      </v-toolbar-title>
    </v-toolbar>
    
    
    <v-content>
      <v-container align-start>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  //import { mapState, mapGetters, mapActions } from 'vuex'
  
  export default {
    data: () => ({
      drawer: null,
      items: [
        { icon: 'camera', text: 'Cameras', route: '/' },
      ],
      isLoggedIn: true // hax for now
    }),
    computed: {
      
    },
    props: {
      source: String
    },
    methods: {
      
    }
  }
</script>
  
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
</style>