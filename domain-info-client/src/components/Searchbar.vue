<template>
  <div class="searchbar-component">
    <b-row>
      <b-col cols="8">
        <b-form-input class="searchbar" v-model="text" placeholder="Enter domain name"></b-form-input>
        
      </b-col>
      <b-col cols="4">
        <b-button class="searchbar-btn" variant="primary">Search</b-button>
      </b-col>
    </b-row>
  </div>
  
</template>

<script>
export default {
  data() {
    return {
      text: ''
    }
  },
  methods: {
    getMeal() {
      
    },
  }
}    
</script>

<style>

.searchbar, .searchbar-btn {
  width: 100%;
} 

</style>