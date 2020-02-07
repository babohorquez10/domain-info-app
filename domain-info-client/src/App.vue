<template>
  <div id="app">
    <b-container>
      <b-row>
        <b-col cols="6">
          <Searchbar />
          <SearchHistory :items=history />

        </b-col>
        <b-col cols="6">
          <ServerInfo :server=server />
        </b-col>
        
      </b-row>
    </b-container>
  </div>
</template>

<script>
import Searchbar from './components/Searchbar.vue'
import ServerInfo from './components/ServerInfo.vue'
import SearchHistory from './components/SearchHistory.vue'
import axios from 'axios'

export default {
  name: 'app',
  components: {
    Searchbar,
    SearchHistory,
    ServerInfo
  },
  data() {
    return {
      server: {
        name: 'Propiedad 1', 
        sslGrade: 'A'
      },
      history: [],
    }
  },
  mounted () {
    axios
      .get('http://localhost:8081/history')
      .then(response => (this.history = response.data.items))
  }
}
</script>

<style>
/* #app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
} */
</style>
