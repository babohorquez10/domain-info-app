<template>
  
  <div class="search-history-component">
    <SearchHistoryItem v-for="item in items" v-bind:key="item.id" :item=item />
  </div>

</template>

<script>

import SearchHistoryItem from './SearchHistoryItem.vue'

export default {
  props: ['items'],
  components : {
    SearchHistoryItem
  },
}    
</script>