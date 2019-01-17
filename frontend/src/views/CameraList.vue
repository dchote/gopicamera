<template>
 <v-layout>
    <v-flex >
      <v-card>
        <v-container fluid grid-list-md>
          <v-layout row wrap>
            <v-flex xs12 v-for="camera in cameras" :key="camera.name">
              <v-card flat tile>
                <v-img :src="camera.cameraURL" class="white--text">
                  <v-container fill-height fluid>
                    <v-layout fill-height>
                      <v-flex xs12 align-end flexbox>
                        <span class="headline">{{ camera.name }}</span>
                      </v-flex>
                    </v-layout>
                  </v-container>
                </v-img>
              </v-card>
            </v-flex>
          </v-layout>
        </v-container>
        </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
  import { mapState } from 'vuex'
  
  export default {
    data: () => ({
      
    }),
    computed: {
      ...mapState('cameras', ['cameras']),
    },
    methods: {
      
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>