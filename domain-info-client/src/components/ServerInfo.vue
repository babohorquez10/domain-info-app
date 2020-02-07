<template>
  
  <div class="server-info-component">
    <p>
      {{ server.name }}
    </p>
    <p>
      {{ server.sslGrade }}
    </p>
  </div>
  
</template>

<script>

export default {
  props: ['server'],
}    
</script>